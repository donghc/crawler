package main

import (
	"fmt"
	"github.com/donghc/crawler/engine"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/log"
	"github.com/donghc/crawler/parse/doubangroup"
	"github.com/donghc/crawler/proxy"
)

func main() {
	doubanGroup()
}

func doubanGroup() {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)

	cookie := "bid=qZ8gG_-P5M0; _pk_id.100001.8cb4=a4aa81c52ee51832.1691651193.; __yadk_uid=gfuYKDTSWELTL3H7yJ78VKbqNVJFAQeG; ll=\"108288\"; __utmz=30149280.1696663356.5.5.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; viewed=\"1007305\"; ap_v=0,6.0; __utmc=30149280; push_noty_num=0; push_doumail_num=0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1700120445%2C%22https%3A%2F%2Ftime.geekbang.com%2Fcolumn%2Farticle%2F612328%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.1471592734.1691651195.1700116287.1700120447.9; __utmt=1; __utmt_douban=1; loc-last-index-location-id=\"108288\"; _vwo_uuid_v2=DCBB38EF9DF817F80715096C3D1EE293B|d2f167ec2e012b3f6a02c950e5ce9b4d; dbcl2=\"167037167:j9y6vSfLpN8\"; ck=iZmW; __utmv=30149280.16703; __utmb=30149280.19.3.1700120488122"

	var seeds []*collect.Request
	str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", 0)
	seeds = append(seeds, &collect.Request{
		URL:       str,
		Cookie:    cookie,
		ParseFunc: doubangroup.ParseURL,
	})
	// for i := 0; i <= 100; i += 25 {
	// 	str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
	// 	workList = append(workList, &collect.Request{
	// 		URL:       str,
	// 		Cookie:    cookie,
	// 		ParseFunc: doubangroup.ParseURL,
	// 	})
	// }
	p, _ := getProxy()

	f := &collect.BrowserFetch{
		Timeout: 3 * time.Second,
		Proxy:   p,
	}

	schedule := engine.NewSchedule(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
	)

	schedule.Run()

}

func getProxy() (proxy.ProxyFunc, error) {
	return nil, nil
	// proxy
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	return proxy.RoundRobinProxySwitcher(proxyURLs...)

}

func test(workList []*collect.Request, f *collect.BrowserFetch) {

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			body, err := f.Get(item)
			if err != nil {
				logger.Error("read content failed ", zap.Error(err))
				continue
			}
			res := item.ParseFunc(body, item)
			for _, v := range res.Items {
				logger.Info("result", zap.String("get url:", v.(string)))
			}
			time.Sleep(5 * time.Second)
			workList = append(workList, res.Requests...)
		}
	}
}
