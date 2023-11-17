package engine

import (
	"github.com/donghc/crawler/collect"
	"go.uber.org/zap"
)

type Schedule struct {
	requestCh chan *collect.Request // requestCh 通道接收来自外界的请求 并将请求存储到 reqQueue 队列中
	workerCh  chan *collect.Request //
	out       chan collect.ParseResult
	options
}

func NewSchedule(opts ...Option) *Schedule {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	s := &Schedule{}
	s.options = o
	return s

}

func (s *Schedule) Run() {
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)

	s.requestCh = requestCh
	s.workerCh = workerCh
	s.out = out

	go s.Schedule()

	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWorker()
	}

	s.HandleResult()
}

func (s *Schedule) Schedule() {
	var reqQueue = s.Seeds
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request
			//如果任务队列 reqQueue 大于 0，意味着有爬虫任务，这时我们获取队列中第一个任务，并将其剔除出队列
			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workerCh
			}
			select {
			case r := <-s.requestCh:
				reqQueue = append(reqQueue, r)
			case ch <- req:
			}
		}
	}()
}

// CreateWorker 创建 worker
func (s *Schedule) CreateWorker() {
	for {
		r := <-s.workerCh
		s.Logger.Sugar().Info("begin get url ", r.URL)
		body, err := s.Fetch.Get(r)
		if err != nil {
			s.Logger.Error("can not fetch ", zap.Error(err))
			continue
		}
		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

// HandleResult 处理爬取并解析后的数据结构
func (s *Schedule) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				// 继续添加
				s.Logger.Sugar().Info("继续添加爬虫链接 ", req.URL)
				s.requestCh <- req
			}
			for _, item := range result.Items {
				s.Logger.Sugar().Info("获取结果 ", item)
			}

		}
	}
}
