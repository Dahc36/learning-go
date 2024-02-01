package identity

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// We use a non-exported type for the keys in context.Context because the equality
// check for ctx.Value(key) checks for the keys' type, so there can be no conflict with
// a different package adding keys to the context.Context
type key int

const (
	userKey key = iota
)

// The exported functions should provide access to reading and writing the actual values
func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userKey).(string)
	return user, ok
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := uuid.New().String()
		if hUser := r.Header.Get("X-User"); hUser != "" {
			user = hUser
		}
		ctx = context.WithValue(ctx, userKey, user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
