package main

import (
	"time"

	"github.com/donghc/crawler/collect/impl"
	"github.com/donghc/crawler/engine"
	"github.com/donghc/crawler/storage/sqlstorage"

	"go.uber.org/zap/zapcore"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/log"
	"github.com/donghc/crawler/proxy"
)

var (
	cookie = "bid=qZ8gG_-P5M0; ll=\"108288\"; __utmz=30149280.1696663356.5.5.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; viewed=\"1007305\"; _pk_id.100001.3ac3=3b24b547be1b35c6.1699849766.; push_noty_num=0; push_doumail_num=0; _vwo_uuid_v2=DCBB38EF9DF817F80715096C3D1EE293B|d2f167ec2e012b3f6a02c950e5ce9b4d; __utmv=30149280.16703; ct=y; __utmz=81379588.1700806485.4.3.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/misc/sorry; _vwo_uuid_v2=DCBB38EF9DF817F80715096C3D1EE293B|d2f167ec2e012b3f6a02c950e5ce9b4d; dbcl2=\"167037167:gpb9vz+25eM\"; ck=Qbp3; __utmc=30149280; __utmc=81379588; __yadk_uid=eyZktnpcn7iCE2bDNLEIRVooJh4Ls6PF; frodotk_db=\"09248aa50ae8fa331b6d9fef20baee99\"; ap_v=0,6.0; _pk_ref.100001.3ac3=%5B%22%22%2C%22%22%2C1701852235%2C%22https%3A%2F%2Fwww.douban.com%2Fmisc%2Fsorry%3Foriginal-url%3Dhttps%3A%2F%2Fbook.douban.com%2Fsubject%2F1007305%2F%22%5D; _pk_ses.100001.3ac3=1; __utma=30149280.1471592734.1691651195.1701842811.1701852235.16; __utma=81379588.1587388512.1699849766.1701842811.1701852235.6; __utmt_douban=1; __utmb=30149280.5.10.1701852235; __utmt=1; __utmb=81379588.5.10.1701852235"
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

	var seeds = make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			// Name: "find_douban_sun_room",
			// Name: "js_find_douban_sun_room",
			Name: "douban_book_list",
		},
		Fetcher: f,
		Storage: storage,
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
