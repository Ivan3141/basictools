package concurrent_middleware

import "context"

type Handler func() error
type PoolMiddleware func(ctx context.Context, handler Handler) Handler

func ChainMiddleware(poolMiddleware ...PoolMiddleware) PoolMiddleware {
	n := len(poolMiddleware)
	return func(ctx context.Context, handler Handler) Handler {
		chainer := func(poolMiddleware PoolMiddleware, currentHandler Handler) Handler {
			return poolMiddleware(ctx, currentHandler)
		}
		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(poolMiddleware[i], chainedHandler)
		}
		return chainedHandler
	}
}
