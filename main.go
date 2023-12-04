package main

import (
	"time"

	"github.com/donghc/crawler/collect/impl"

	"github.com/donghc/crawler/engine"

	"go.uber.org/zap/zapcore"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/log"
	"github.com/donghc/crawler/proxy"
)

var (
	cookie = "bid=tC5jgShcUIU; ll=\"108288\"; _pk_id.100001.8cb4=baf74bc962a73268.1700147162.; __utmc=30149280; douban-fav-remind=1; __yadk_uid=Lob2AXbySoLVi5shgZY82uV8bFP4TSd3; ct=y; push_doumail_num=0; __utmz=30149280.1700924224.9.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.16703; push_noty_num=0; frodotk_db=\"f912931b88658561eb37b95299309023\"; ps=y; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1701698994%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.1262193242.1695303769.1701694458.1701699006.13; dbcl2=\"167037167:gpb9vz+25eM\"; ck=Qbp3; __utmt=1; __utmb=30149280.34.5.1701702326167"
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

	var seeds = make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			Name: "js_find_douban_sun_room",
		},
		Fetcher: f,
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
