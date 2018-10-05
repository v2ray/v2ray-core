package dispatcher

import "context"

//go:generate errorgen

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
