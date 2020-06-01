package aratest

import (
	"context"
)

type cxtKey int

const (
	hostPortKey cxtKey = iota
	tokenKey
)

func withHostPort(ctx context.Context, hostport string) context.Context {
	return context.WithValue(ctx, hostPortKey, hostport)
}

func hostPortFromCtx(ctx context.Context) string {
	return ctx.Value(hostPortKey).(string)
}

func withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func tokenFromCtx(ctx context.Context) string {
	var token, ok = ctx.Value(tokenKey).(string)
	if !ok {
		return ""
	}
	return token
}
