package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
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

type AuthController struct {
	auth  *internal.AuthService
	rdb   redis.UniversalClient
	users *internal.UserService
}

func NewAuthController(rdb redis.UniversalClient, auth *internal.AuthService, users *internal.UserService) *AuthController {
	return &AuthController{
		auth:  auth,
		rdb:   rdb,
		users: users,
	}
}

var ErrUsernameTaken = errors.New("username is taken")
var ErrInvalidUsername = errors.New("username is invalid")

func useCookieAuth(req *http.Request) bool {
	return !strings.Contains(req.UserAgent(), "pizzatribes")
}

func (a *AuthController) Login(ctx context.Context, username, password string) (string, error) {
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

func (a *AuthController) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		req := registerRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = a.users.CreateUser(r.Context(), req.Username, req.Password)
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
		jwt, err := a.auth.CreateToken(userId, expiresAt)
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
