package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	etcdReg "github.com/go-micro/plugins/v4/registry/etcd"
	gs "github.com/go-micro/plugins/v4/server/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"github.com/donghc/crawler/collect/impl"
	"github.com/donghc/crawler/engine"
	"github.com/donghc/crawler/limiter"
	"github.com/donghc/crawler/proto/greeter"
	"github.com/donghc/crawler/storage/sqlstorage"

	"go.uber.org/zap/zapcore"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/log"
	"github.com/donghc/crawler/proxy"
)

var (
	plugin = log.NewStdoutPlugin(zapcore.InfoLevel)
	logger = log.NewLogger(plugin)
)

func main() {
	// go doubanGroup()
	go HandleHTTP()
	register()
}

func doubanGroup() {

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

func register() {

	newRegistry := etcdReg.NewRegistry(registry.Addrs(":2379"))

	service := micro.NewService(
		micro.Server(gs.NewServer(server.Id("worker1"))),
		micro.Address(":9090"),
		micro.Name("go.micro.server.worker"),
		micro.Registry(newRegistry),
		micro.WrapHandler(logWrapper(logger)),
	)

	service.Init()

	greeter.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		logger.Fatal("grpc server stop")
	}
}

func HandleHTTP() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := greeter.RegisterGreeterGwFromEndpoint(ctx, mux, ":9090", opts)
	if err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(":8080", mux)
}

type Greeter struct {
}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, resp *greeter.Response) (err error) {
	resp.Greeting = "hello" + req.Name

	return nil
}

func logWrapper(log *zap.Logger) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Info("recieve request",
				zap.String("method", req.Method()),
				zap.String("Service", req.Service()),
				zap.Reflect("request param:", req.Body()),
			)
			err := fn(ctx, req, rsp)
			return err
		}
	}
}
