package aratest

import (
	"context"
	"net/url"
)

type cxtKey int

const (
	hostPortKey cxtKey = iota
	tokenKey
)

func withHostPort(ctx context.Context, hostport url.URL) context.Context {
	return context.WithValue(ctx, hostPortKey, hostport)
}

func hostPortFromCtx(ctx context.Context) url.URL {
	return ctx.Value(hostPortKey)
}

func withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func tokenFromCtx(ctx context.Context) string {
	return ctx.Value(tokenKey)
}
