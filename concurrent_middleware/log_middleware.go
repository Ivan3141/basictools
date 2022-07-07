package concurrent_middleware

import (
	"git.bybit.com/svc/go/pkg/bstd"
	"git.bybit.com/svc/go/pkg/bzap"
	"git.bybit.com/yan.fan/concurrent-package-ian/codes"
	"go.uber.org/zap"
)

func LogMiddleWare(parentLogger *zap.Logger, handler Handler) Handler {
	return func() (err error) {
		var logger *zap.Logger

		logger = parentLogger.With(zap.String("method", ""))

		defer func() {
			err = handleError(logger, err.(*codes.Code))
		}()

		err = bstd.RunWithRecover(func() error {
			err = handler()
			return err
		})

		return err
	}
}

func handleError(logger *zap.Logger, err *codes.Code) error {
	if err == nil {
		return nil
	}
	bzap.LogError(logger, "run_error", err)
	return err
}
