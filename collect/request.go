package collect

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/donghc/crawler/storage"
)

type Property struct {
	Name     string        `json:"name,omitempty"` // 任务名称，需要保证唯一
	URL      string        `json:"url,omitempty"`  // URL地址
	Cookie   string        `json:"cookie"`         // cookie
	WaitTime time.Duration `json:"wait_time"`      // 等待时间
	Reload   bool          `json:"reload"`         // 是否可以重复爬取
	MaxDepth int64         `json:"max_depth"`      // 最大深度
}

// Task 爬虫任务
type Task struct {
	Property

	RootReq     *Request // 任务中的第一请求
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher

	Rule RuleTree //

	Storage storage.Storage
	Logger  *zap.Logger
}

// RuleContext 为自定义结构体，用于传递上下文信息，也就是当前的请求参数以及要解析的内容字节数组。
type RuleContext struct {
	Body []byte
	Req  *Request
}

func (c *RuleContext) ParseJSReg(name string, reg string) ParseResult {
	r2 := regexp.MustCompile("禁止访问")
	ok := r2.MatchString(string(c.Body))
	if ok {
		fmt.Println("被检测到为爬虫行为，需要人机验证")
		return ParseResult{}
	}

	re := regexp.MustCompile(reg)
	matches := re.FindAllSubmatch(c.Body, -1)
	result := ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &Request{
			Method:   http.MethodGet,
			Task:     c.Req.Task,
			URL:      u,
			Depth:    c.Req.Depth + 1,
			RuleName: name,
		})
	}
	return result
}

func (c *RuleContext) OutputJS(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	resultStr := re.FindString(string(c.Body))

	r2 := regexp.MustCompile("阳台")
	ok := r2.MatchString(resultStr)
	if !ok {
		return ParseResult{
			Items: []interface{}{},
		}
	}

	result := ParseResult{
		Items: []interface{}{c.Req.URL},
	}

	return result
}

func (c *RuleContext) Output(data interface{}) *storage.DataCell {
	res := &storage.DataCell{}
	res.Data = make(map[string]interface{})
	res.Data["Task"] = c.Req.Task.Name
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["URL"] = c.Req.URL
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")

	return res
}

type Request struct {
	unique   string // 唯一标识
	RuleName string // 规则名称
	URL      string // 表示要访问的网站
	Method   string // 方法
	Depth    int64  // 当前深度
	Priority int64  // 优先级

	Task *Task

	TmpData *Temp
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
