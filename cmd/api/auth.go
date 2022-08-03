package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

/*

This is an implementation of a simple user system (with login/logout/register).
- Users are stored in Redis
- Authentication is done using JWTs
- Passwords are hashed using bcrypt


Redis details

- Users are stored as hashes in user:{userid}
- User ids can be looked up using username:{username}

*/

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"accessToken"`
}

type userDbo struct {
	Id             string `redis:"id"`
	Username       string `redis:"username"`
	HashedPassword string `redis:"hashed_password"`
}

type AuthService struct {
	internal.AuthService
	rdb redis.UniversalClient
}

func NewAuthService(rdb redis.UniversalClient) *AuthService {
	return &AuthService{
		rdb: rdb,
	}
}

var ErrUsernameTaken = errors.New("username is taken")
var ErrInvalidUsername = errors.New("username is invalid")

func useCookieAuth(req *http.Request) bool {
	return !strings.Contains(req.UserAgent(), "pizzatribes")
}

var validUsername = regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z0-9_]*$`)

func IsValidUsername(username string) bool {
	return validUsername.MatchString(username) && len(username) >= 3 && len(username) <= 20
}

func (a *AuthService) Register(ctx context.Context, username, password string) error {
	if !IsValidUsername(username) {
		return ErrInvalidUsername
	}

	id := xid.New().String()
	usernameKey := fmt.Sprintf("username:%s", strings.ToLower(username))
	userKey := fmt.Sprintf("user:%s", id)

	// Check for existing user with this username
	res, err := a.rdb.Exists(ctx, usernameKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if res != 0 {
		return ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

			pipe.Set(ctx, usernameKey, id, 0)
			pipe.HSet(ctx, userKey, "id", id, "username", username, "hashed_password", hash)
			pipe.SAdd(ctx, "users", id)
			return nil
		})
		return err
	}

	err = a.rdb.Watch(ctx, txf, usernameKey, userKey)

	return err
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	usernameKey := fmt.Sprintf("username:%s", strings.ToLower(username))
	userId, err := a.rdb.Get(ctx, usernameKey).Result()
	if err != nil {
		return "", err
	}

	userKey := fmt.Sprintf("user:%s", userId)
	user := userDbo{}
	err = a.rdb.HGetAll(ctx, userKey).Scan(&user)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func (a *AuthService) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		req := registerRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = a.Register(r.Context(), req.Username, req.Password)
		if err != nil {
			log.Error().Err(err).Msg("Error when registering user")
			if errors.Is(err, ErrUsernameTaken) {
				http.Error(w, "Username taken", http.StatusBadRequest)
			} else if errors.Is(err, ErrInvalidUsername) {
				http.Error(w, "Invalid username", http.StatusBadRequest)
			} else {
				http.Error(w, "Failed to register", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(200)
		w.Write([]byte("{}"))
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		req := loginRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse login request")
			http.Error(w, "Invalid login request", http.StatusBadRequest)
			return
		}

		userId, err := a.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			if err == redis.Nil {
				log.Info().Msg("Bad credentials: no such user")
				http.Error(w, "Bad credentials", http.StatusForbidden)
				return
			}

			if err == bcrypt.ErrMismatchedHashAndPassword {
				log.Info().Msg("Bad credentials: password mismatch")
				http.Error(w, "Bad credentials", http.StatusForbidden)
				return
			}

			log.Error().Err(err).Msg("Login failed")
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}
		expiresIn := 3 * 7 * 24 * time.Hour // 3 weeks
		expiresAt := time.Now().Add(expiresIn)
		jwt, err := a.CreateToken(userId, expiresAt)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			return
		}

		if useCookieAuth(r) {
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    jwt,
				HttpOnly: true,
				Path:     "/",
				MaxAge:   int(expiresIn.Seconds()),
				SameSite: http.SameSiteStrictMode,
			})
			w.WriteHeader(200)
			w.Write([]byte("{}"))
		} else {
			b, err := json.Marshal(&loginResponse{
				AccessToken: jwt,
			})
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal login response")
				w.WriteHeader(500)
				return
			}

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(b)
		}
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(200)
		w.Write([]byte("{}"))
	})

	return r
}
