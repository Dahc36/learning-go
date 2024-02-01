package tracker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dahc36/learning-go/12-the-context/pkg/identity"
	"github.com/google/uuid"
)

type key int

const (
	guidKey key = iota
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, guidKey, uuid.New().String())
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

type Logger struct{}

func (Logger) Log(ctx context.Context, message string) {
	user, ok := identity.UserFromContext(ctx)
	if ok {
		message = fmt.Sprintf("User: %s - %s", user, message)
	}
	guid, ok := ctx.Value(guidKey).(string)
	if ok {
		message = fmt.Sprintf("GUID: %s - %s", guid, message)
	}
	fmt.Println(message)
}
