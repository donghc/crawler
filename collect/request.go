package collect

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

type Property struct {
	Name     string        `json:"name,omitempty"` // 任务名称，需要保证唯一
	URL      string        `json:"url,omitempty"`  // URL地址
	Cookie   string        `json:"cookie"`         // cookie
	WaitTime time.Duration `json:"wait_time"`      // 等待时间
	Reload   bool          `json:"reload"`         // 是否可以重复爬取
	MaxDepth int64         `json:"max_depth"`      //最大深度
}

// Task 爬虫任务
type Task struct {
	Property

	RootReq     *Request // 任务中的第一请求
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher

	Rule RuleTree //
}

type Request struct {
	unique   string // 唯一标识
	RuleName string //规则名称
	URL      string // 表示要访问的网站
	Method   string // 方法
	Depth    int64  // 当前深度
	Priority int    // 优先级

	Task *Task
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
