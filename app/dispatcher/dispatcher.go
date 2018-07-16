package dispatcher

import "context"

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dispatcher -path App,Dispatcher

type key int

const (
	sniffing key = iota
)

func ContextWithSniffingResult(ctx context.Context, r SniffResult) context.Context {
	return context.WithValue(ctx, sniffing, r)
}

func SniffingResultFromContext(ctx context.Context) SniffResult {
	if c, ok := ctx.Value(sniffing).(SniffResult); ok {
		return c
	}
	return nil
}
