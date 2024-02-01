package tracker

import "context"

// We use a non-exported type for the keys in context.Context because the equality
// check for ctx.Value(key) checks for the keys' type, so there can be no conflict with
// a different package adding keys to the context.Context
type key int

const (
	guidKey key = iota
	otherKey
)

// The exported functions should provide access to reading and writing the actual values
func ContextWithGUID(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, guidKey, user)
}
func UerFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(guidKey).(string)
	return user, ok
}
