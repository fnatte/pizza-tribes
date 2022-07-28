package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

type AuthService struct {}

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

func (a *AuthService) CreateToken(userId string, expiresAt time.Time) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Subject:   userId,
	}

	tokenString, err := t.SignedString(jwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
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

