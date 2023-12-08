package doubangroup

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/donghc/crawler/collect"
)

var (
	cookie = "bid=tC5jgShcUIU; ll=\"108288\"; __utmc=30149280; douban-fav-remind=1; ct=y; ap_v=0,6.0; dbcl2=\"167037167:j9y6vSfLpN8\"; ck=iZmW; push_doumail_num=0; __utma=30149280.1262193242.1695303769.1700921222.1700924224.9; __utmz=30149280.1700924224.9.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmv=30149280.16703; push_noty_num=0; __utmb=30149280.11.5.1700924241581; frodotk_db=\"f912931b88658561eb37b95299309023\"; ps=y"
)

var DoubanGroupTask = &collect.Task{
	Property: collect.Property{
		Name:     "find_douban_sun_room",
		Cookie:   cookie,
		WaitTime: 3,
		MaxDepth: 5,
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			var roots []*collect.Request
			for i := 0; i < 25; i += 25 {
				str := fmt.Sprintf("https://www.douban.com/group/beijingzufang/discussion?start=%d&type=new", i)
				roots = append(roots, &collect.Request{
					Priority: 1,
					URL:      str,
					Method:   http.MethodGet,
					RuleName: "解析网站URL",
				},
				)
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"解析网站URL": &collect.Rule{ParseFunc: ParseURL},
			"解析阳台房":   &collect.Rule{ParseFunc: GetSunRoom},
		},
	},
}

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(ctx *collect.RuleContext) (collect.ParseResult, error) {
	r2 := regexp.MustCompile("禁止访问")
	ok := r2.MatchString(string(ctx.Body))
	if ok {
		fmt.Println("被检测到为爬虫行为，需要人机验证")
		return collect.ParseResult{}, errors.New("被检测到为爬虫行为，需要人机验证")
	}

	re := regexp.MustCompile(cityListRe)
	result := collect.ParseResult{}
	matches := re.FindAllSubmatch(ctx.Body, -1)

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Method:   http.MethodGet,
			Task:     ctx.Req.Task,
			URL:      u,
			Depth:    ctx.Req.Depth + 1,
			RuleName: "解析阳台房",
		})
	}
	return result, nil
}

const ContentRe = `<div\s+class="topic-content">(?s:.)*?</div>`

func GetSunRoom(ctx *collect.RuleContext) (collect.ParseResult, error) {
	re := regexp.MustCompile(ContentRe)
	resultStr := re.FindString(string(ctx.Body))

	r2 := regexp.MustCompile("阳台")

	ok := r2.MatchString(resultStr)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}, nil
	}

	result := collect.ParseResult{
		Items: []interface{}{ctx.Req.URL},
	}

	return result, nil
}
