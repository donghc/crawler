package engine

import (
	"github.com/donghc/crawler/collect"
	"go.uber.org/zap"
)

type Schedule struct {
	requestCh chan *collect.Request // requestCh 通道接收来自外界的请求 并将请求存储到 reqQueue 队列中
	workerCh  chan *collect.Request //
	reqQueue  []*collect.Request    //
	Logger    *zap.Logger
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	s.requestCh = make(chan *collect.Request)
	s.workerCh = make(chan *collect.Request)

	return s
}

func (s *Schedule) Push(reqs ...*collect.Request) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *collect.Request {
	r := <-s.workerCh
	return r
}

func (s *Schedule) Schedule() {
	var req *collect.Request
	var ch chan *collect.Request
	for {
		//如果任务队列 reqQueue 大于 0，意味着有爬虫任务，这时我们获取队列中第一个任务，并将其剔除出队列
		if len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}
		select {
		case r := <-s.requestCh:
			if req != nil {
				s.reqQueue = append(s.reqQueue, req)
			}
			s.reqQueue = append(s.reqQueue, r)
		case ch <- req:
		}
	}
}
