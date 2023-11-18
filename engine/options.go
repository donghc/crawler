package engine

import (
	"github.com/donghc/crawler/collect"
	"go.uber.org/zap"
)

type Option func(opts *options)

type options struct {
	Seeds     []*collect.Task
	WorkCount int //并发度
	Fetch     collect.Fetcher
	Logger    *zap.Logger
}

var defaultOptions = options{
	Logger: zap.NewNop(),
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}

func WithFetcher(f collect.Fetcher) Option {
	return func(opts *options) {
		opts.Fetch = f
	}
}
func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}
func WithSeeds(seed []*collect.Task) Option {
	return func(opts *options) {
		opts.Seeds = seed
	}
}
