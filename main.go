package main

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/donghc/crawler/collect/impl"
	"github.com/donghc/crawler/engine"
	"github.com/donghc/crawler/limiter"
	"github.com/donghc/crawler/storage/sqlstorage"

	"go.uber.org/zap/zapcore"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/log"
	"github.com/donghc/crawler/proxy"
)

func main() {
	doubanGroup()
}

func doubanGroup() {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)

	logger.Info("log init end ,begin start crawler task")

	p, _ := getProxy()
	f := &impl.BrowserFetch{
		Timeout: 3 * time.Second,
		Proxy:   p,
	}

	storage, err := sqlstorage.NewSqlStore(
		sqlstorage.WithLogger(logger.Named("sqldb")),
		sqlstorage.WithBatchCount(2),
		sqlstorage.WithSqlUrl("root:123456@tcp(127.0.0.1:3306)/lsb?charset=utf8"),
	)
	if err != nil {
		logger.Panic("create sql storage failed")
		return
	}
	// 表示每3秒1个令牌
	secondLimit := rate.NewLimiter(limiter.Per(1, 3*time.Second), 1)
	// 60秒20个
	minuteLimit := rate.NewLimiter(limiter.Per(20, 1*time.Minute), 20)

	multiLimiter := limiter.NewMultiLimiter(secondLimit, minuteLimit)

	var seeds = make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			// Name: "find_douban_sun_room",
			// Name: "js_find_douban_sun_room",
			Name: "douban_book_list",
		},
		Fetcher: f,
		Storage: storage,
		Limiter: multiLimiter,
	},
	)

	schedule := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(1),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	schedule.Run()

}

func getProxy() (proxy.ProxyFunc, error) {
	return nil, nil
	// proxy
	proxyURLs := []string{"http://58.246.58.150:9002"}
	return proxy.RoundRobinProxySwitcher(proxyURLs...)

}
