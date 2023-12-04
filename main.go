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
	cookie = "BDUSS_BFESS=9FU1BPQ0EwOWZPbmE3RHFHNEJLNWdhYnN1WnN1YzR2LTU5cWl6SGROZXAxfkZrRUFBQUFBJCQAAAAAAAAAAAEAAAAt1ekHgs7RYrXE0LDEpwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKlKymSpSspkcW; H_WISE_SIDS_BFESS=234020_131861_114551_216851_213357_214790_110085_244722_254835_261715_236312_259308_263346_256419_256223_265881_266362_266844_265615_265985_256302_266188_267898_265776_264354_266713_268592_268689_266186_268849_259642_269236_269731_269776_268235_269904_270084_267066_256739_270460_267528_270598_270548_270834_271113_271120_271170_271178_271074_271193_268727_268987_269034_271225_267659_271322_265032_269609_271471_266028_271581_270482_270102_271562_271867_270157_269771_271812_269878_271934_269563_267807_269664_256151_269212_234296_234208_272082_266324_272163_270181_270054_272280_263618_266566_272335_272364_272456_271144_272511_253022_271545_271905_271689_270186_272817_272822_272641_272839_272801_270142_260335_269297_269715_267596_273059_273092_203519_273165_273146_273212_273237_273301_273399; BAIDUID_BFESS=1B5EE11CA94CBB4CEB10E8A5DCAED18A:SL=0:NR=10:FG=1; ZFY=KHVfhTx:B2BJyZ:BZwf1Ea8GVb67mlSMlVBhtnzfzn2pw:C; ab_sr=1.0.1_YTk3NTRiZWFkNmMwMDIxZDYwMDg4NGJkZGIyMDgzNDc3Nzg2ZmI2MzZjOWFhZTU0ZDE4MGI2ZDI4ZmQ4NjdhMjhiYmQzMTQwMTZkOWExOTQ0ZTYxODRlMjJiYTcyMjgxYWFhZmYyOWJmMDRlYTU5Yjk3MDJkOGY0ZWY2MGVjZWMwODZmNzg5NWEwOGUyMDEwZWJjOWYzNmRjNDAxMjc0MTQwMWYyNWEwM2UxZTE4OGUzYjEyNGYyZmNkMDI4ZTY5"
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
	proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	return proxy.RoundRobinProxySwitcher(proxyURLs...)

}
