package engine

import (
	"go.uber.org/zap"

	"github.com/donghc/crawler/collect"
)

type Schedule struct {
	requestCh   chan *collect.Request // requestCh 通道接收来自外界的请求 并将请求存储到 reqQueue 队列中
	workerCh    chan *collect.Request //
	priReqQueue []*collect.Request    // 优先级队列
	reqQueue    []*collect.Request    //
	Logger      *zap.Logger
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
func (s *Schedule) Output() *collect.Request {
	r := <-s.workerCh
	return r
}

func (s *Schedule) Schedule() {
	var req *collect.Request
	var ch chan *collect.Request
	for {
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerCh
		}
		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}

		select {
		case r := <-s.requestCh:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}
