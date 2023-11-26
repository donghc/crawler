package collect

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type RuleTree struct {
	// RuleTree.Root 是一个函数，用于生成爬虫的种子网站
	Root func() ([]*Request, error)
	// 规则哈希表，用于存储当前任务所有的规则，哈希表的 Key 为规则名，Value 为具体的规则,每一个规则就是一个 ParseFunc 解析函数。
	// 参数 RuleContext 为自定义结构体，用于传递上下文信息，也就是当前的请求参数以及要解析的内容字节数组。
	// 后续还会添加请求中的临时数据等上下文数据
	Trunk map[string]*Rule
}

// Rule 采集规则节点
type Rule struct {
	ParseFunc func(ctx *RuleContext) (ParseResult, error) // 内容解析函数
}

// RuleContext 为自定义结构体，用于传递上下文信息，也就是当前的请求参数以及要解析的内容字节数组。
type RuleContext struct {
	Body []byte
	Req  *Request
}

func (c *RuleContext) ParseJSReg(name string, reg string) (ParseResult, error) {
	r2 := regexp.MustCompile("禁止访问")
	ok := r2.MatchString(string(c.Body))
	if ok {
		fmt.Println("被检测到为爬虫行为，需要人机验证")
		return ParseResult{}, errors.New("被检测到为爬虫行为，需要人机验证")
	}

	re := regexp.MustCompile(reg)
	result := ParseResult{}
	matches := re.FindAllSubmatch(c.Body, -1)

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
	return result, nil
}

func (c *RuleContext) OutputJS(reg string) (ParseResult, error) {
	re := regexp.MustCompile(reg)
	resultStr := re.FindString(string(c.Body))

	r2 := regexp.MustCompile("阳台")
	ok := r2.MatchString(resultStr)
	if !ok {
		return ParseResult{
			Items: []interface{}{},
		}, nil
	}

	result := ParseResult{
		Items: []interface{}{c.Req.URL},
	}

	return result, nil
}
