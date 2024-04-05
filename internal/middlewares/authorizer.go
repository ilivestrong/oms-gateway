package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilivestrong/oms-gateway/internal/auth"
)

const (
	ErrAuthHeaderMissing       = "Authorization header is missing"
	ErrInvalidToken            = "invalid token"
	ErrUnexpectedSigningMethod = "unexpected signing method"
	BearerAuth                 = "Bearer "
	JWTEncryptionAlgoHeader    = "alg"
	AuthorizationHeader        = "Authorization"
	LoginEndpointURL           = "/login"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == LoginEndpointURL {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			http.Error(w, ErrAuthHeaderMissing, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, BearerAuth, "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%s: %v", ErrUnexpectedSigningMethod, token.Header[JWTEncryptionAlgoHeader])
			}
			return []byte(auth.TOKEN_SECRET), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
