package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
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
	rdb redis.UniversalClient
}

func NewAuthService(rdb redis.UniversalClient) *AuthService {
	return &AuthService{
		rdb: rdb,
	}
}

var ErrUsernameTaken = errors.New("username is taken")

var jwtSigningKey = []byte{}

func init() {
	jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))
	if len(jwtSigningKey) == 0 {
		panic("JWT_SIGNING_KEY was not set")
	}
}

func getJwtSigningKeyFunc(*jwt.Token) (interface{}, error) {
	return jwtSigningKey, nil
}

func useCookieAuth(req *http.Request) bool {
	return !strings.Contains(req.UserAgent(), "pizzatribes")
}

func (a *AuthService) Register(ctx context.Context, username, password string) error {
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

func (a *AuthService) CreateToken(userId string) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
		Subject:   userId,
	}

	tokenString, err := t.SignedString(jwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

// Read access token from Authorization header, cookie, or web socket protocol header:
// - Authorization header is read with auth scheme "Bearer". For example, it will read from "Authorization: Bearer {token}".
// - The cookie is simply named "token". For example, it will read from "Cookie: token={token}"
// - For web sockets, which do not support custom headers, we take horrific usage of the Sec-WebSocket-Protcol header,
//   and will read tokens from protocols that start with "accessToken.". For example: "Sec-WebSocket-Protocol: pizzatribes, accessToken.{token}".
func getAccessToken(r *http.Request) (string, error) {
	// From Authorization header
	authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if authHeader[0] == "Bearer" {
		if len(authHeader) < 2 {
			return "", nil
		}
		return authHeader[1], nil
	}

	// From cookie
	cookie, err := r.Cookie("token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", err
	}
	if err == nil {
		return cookie.Value, nil
	}

	// From Sec-WebSocket-Protocol
	// Yes, this is weird, but a quick workaround to do cross-origin websocket authorization.
	// Some background:
	// - https://stackoverflow.com/q/4361173/86298
	// - https://github.com/whatwg/html/issues/3062
	//
	// Further, note that Sec-WebSocket-Protocol may not include any character:
	// - https://github.com/WebKit/webkit/blob/main/Source/WebCore/Modules/websockets/WebSocket.cpp#L83
	// - https://datatracker.ietf.org/doc/html/draft-ietf-hybi-thewebsocketprotocol-10#section-5.1
	wsProtocolHeaders := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",")
	for _, p := range wsProtocolHeaders {
		const accessTokenPrefix = "accessToken."
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, accessTokenPrefix) {
			return p[len(accessTokenPrefix):], nil
		}
	}

	return "", nil
}

func (a *AuthService) Authorize(r *http.Request) error {
	token, err := getAccessToken(r)
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("not authenticated")
	}

	// Now parse the token
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, getJwtSigningKeyFunc)
	if err != nil {
		return err
	}

	alg := parsedToken.Header["alg"]
	if alg != jwt.SigningMethodHS256.Alg() {
		return fmt.Errorf("error validating token algorithm: %s", alg)
	}

	if !parsedToken.Valid {
		return errors.New("token is invalid")
	}

	claims := parsedToken.Claims.(*jwt.StandardClaims)

	ctx := context.WithValue(r.Context(), "userId", claims.Subject)
	newRequest := r.WithContext(ctx)
	*r = *newRequest

	return nil
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

		jwt, err := a.CreateToken(userId)
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
				MaxAge:   3600 * 72,
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
