package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const Key = ContextKey("user")

type TokenPayload struct {
	ID       int
	Username string
	Role     string
	Exp      int64 // Unix timestamp
}

func MiddlwareJWT(publicKey *rsa.PublicKey) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from header
			var headerToken string
			_, err := fmt.Sscanf(r.Header.Get("Authorization"), "Bearer %s", &headerToken)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Parse token
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(headerToken, claims, func(t *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			// error handling
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// is token valid
			// idgaf what is token.Valid but anyway
			if !token.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// check expiration
			if isExpiredToken(claims) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user := TokenPayload{
				ID:       int(claims["id"].(float64)),
				Username: claims["username"].(string),
				Role:     claims["role"].(string),
				Exp:      int64(claims["exp"].(float64)),
			}

			// set context
			ctx := context.WithValue(r.Context(), Key, user)

			// set context
			r = r.WithContext(ctx)

			// next
			next.ServeHTTP(w, r)
		})
	}
}

func FromContext(ctx context.Context) (TokenPayload, error) {
	payload, ok := ctx.Value(Key).(TokenPayload)
	if !ok {
		return TokenPayload{}, fmt.Errorf("user not found")
	}
	return payload, nil
}

func isExpiredToken(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(int64)
	if !ok {
		return false
	}
	return exp <= time.Now().Unix()
}
