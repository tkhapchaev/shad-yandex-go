//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

type User struct {
	Name  string
	Email string
}

type ctxKey string

var ErrInvalidToken = errors.New("invalid token")

func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(ctxKey("user")).(*User)

	return user, ok
}

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			if authorization == "" {
				http.Error(w, "Can't find authorization header", http.StatusUnauthorized)

				return
			}

			auth := strings.SplitN(authorization, " ", 2)

			if len(auth) != 2 || auth[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)

				return
			}

			token := auth[1]
			user, err := checker.CheckToken(r.Context(), token)

			if err != nil {
				if errors.Is(err, ErrInvalidToken) {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
				} else {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}

				return
			}

			ctx := context.WithValue(r.Context(), ctxKey("user"), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
