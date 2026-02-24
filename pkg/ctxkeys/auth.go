package ctxkeys

import "context"

type ctxAuthKey struct{}

func WithAuthKey(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxAuthKey{}, token)
}

func AuthKey(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxAuthKey{}).(string)
	return v, ok
}
