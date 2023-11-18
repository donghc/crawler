package collect

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

// Task 爬虫任务
type Task struct {
	URL         string        // 表示要访问的网站
	Cookie      string        // cookie
	WaitTime    time.Duration // 默认等待时间
	MaxDepth    int           //最大深度
	RootReq     *Request      //任务中的第一请求
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher
}

type Request struct {
	unique    string //唯一标识
	Task      *Task
	URL       string                             // 表示要访问的网站
	Method    string                             // 方法
	Depth     int                                //当前深度
	ParseFunc func([]byte, *Request) ParseResult // ParseFunc 函数会解析从网站获取到的网站信息，并返回 Requests 数组用于进一步获取数据。而 Items 表示获取到的数据。
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{} // Items 表示获取到的数据。
}

// CheckDepth 查看当前深度是否超过最大深度限制
func (req *Request) CheckDepth() bool {
	if req.Depth > req.Task.MaxDepth {
		return true
	}
	return false
}

// Unique 请求的唯一识别码
func (req *Request) Unique() string {
	block := md5.Sum([]byte(req.URL + req.Method))
	return hex.EncodeToString(block[:])
}
