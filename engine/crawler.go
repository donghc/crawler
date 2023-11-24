package engine

import (
	"sync"

	"go.uber.org/zap"

	"github.com/donghc/crawler/collect"
)

// Crawler 全局爬取实例
type Crawler struct {
	out         chan collect.ParseResult
	Visited     map[string]bool
	VisitedLock sync.Mutex
	failures    map[string]*collect.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex
	options
}

type Scheduler interface {
	Schedule()                // 调度
	Push(...*collect.Request) // 将任务放到调度器中
	Pull() *collect.Request   // 获取任务
}

func NewEngine(opts ...Option) *Crawler {
	option := defaultOptions
	for _, opt := range opts {
		opt(&option)
	}
	c := &Crawler{}
	c.Visited = make(map[string]bool, 100)

	out := make(chan collect.ParseResult)
	c.out = out
	c.options = option

	return c

}

func (c *Crawler) Run() {
	for i := 0; i < c.WorkCount; i++ {
		go c.CreateWorker()
	}

	go c.Schedule()

	c.HandleResult()
}

// CreateWorker 创建 worker
func (c *Crawler) CreateWorker() {
	for {
		r := c.scheduler.Pull()
		if r.CheckDepth() {
			c.Logger.Sugar().Warn("current depth 超过最大限制")
			continue
		}
		if c.HasVisited(r) {
			c.Logger.Debug("此 URL 已经被请求过", zap.String("url :", r.URL))
			continue
		}
		c.StoreVisited(r)

		body, err := r.Task.Fetcher.Get(r)
		if err != nil {
			c.Logger.Error("can not fetch ", zap.Error(err))
			c.SetFailure(r)
			continue
		}
		if len(body) < 6000 {
			c.Logger.Error("can not fetch ", zap.Int("length :", len(body)), zap.String("url :", r.URL))
			c.SetFailure(r)
			continue
		}
		result := r.ParseFunc(body, r)

		if len(result.Requests) > 0 {
			go c.scheduler.Push(result.Requests...)
		}
		c.out <- result
	}
}

func (c *Crawler) Schedule() {
	var reqs []*collect.Request
	for _, seed := range c.Seeds {
		seed.RootReq.Task = seed
		seed.RootReq.URL = seed.URL
		reqs = append(reqs, seed.RootReq)
	}
	go c.scheduler.Schedule()
	go c.scheduler.Push(reqs...)
}

// HandleResult 处理爬取并解析后的数据结构
func (c *Crawler) HandleResult() {
	for {
		select {
		case result := <-c.out:
			for _, item := range result.Items {
				// todo : store
				c.Logger.Sugar().Info("获取结果 ", item)
			}

		}
	}
}

// HasVisited 判断当前请求是否已经被访问过
func (c *Crawler) HasVisited(r *collect.Request) bool {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()
	unique := r.Unique()
	return c.Visited[unique]
}

func (c *Crawler) StoreVisited(reqs ...*collect.Request) {
	c.VisitedLock.Lock()
	defer c.VisitedLock.Unlock()

	for _, req := range reqs {
		unique := req.Unique()
		c.Visited[unique] = true
	}
}

func (c *Crawler) SetFailure(req *collect.Request) {
	unique := req.Unique()
	if !req.Task.Reload {
		// 不可以重复爬取,需要在失败重试前删除 Visited 中的历史记录
		c.VisitedLock.Lock()
		delete(c.Visited, unique)
		c.VisitedLock.Unlock()
		return
	}
	c.failureLock.Lock()
	defer c.failureLock.Unlock()

	if _, ok := c.failures[unique]; !ok {
		// 第一次失败，可以重试一次
		c.failures[unique] = req
		c.scheduler.Push(req)
	}
	// todo: 失败2次，加载到失败队列中
}
