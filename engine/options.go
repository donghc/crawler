package engine

import (
	"go.uber.org/zap"

	"github.com/donghc/crawler/collect"
)

type Option func(opts *options)

type options struct {
	WorkCount int // 并发度
	Fetcher   collect.Fetcher
	Seeds     []*collect.Task
	Logger    *zap.Logger
	scheduler Scheduler
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
		opts.Fetcher = f
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

func WithScheduler(scheduler Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = scheduler
	}
}
