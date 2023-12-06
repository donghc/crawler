package mysql

import (
	"go.uber.org/zap"
)

type options struct {
	logger *zap.Logger
	sqlDSN string
}

var defaultOptions = options{logger: zap.NewNop()}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithConnectionURL(dsn string) Option {
	return func(opts *options) {
		opts.sqlDSN = dsn
	}
}
