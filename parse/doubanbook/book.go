package doubanbook

import (
	"fmt"
	"github.com/donghc/crawler/collect"
	"regexp"
	"strconv"
	"time"
)

var (
	cookie = "bid=tC5jgShcUIU; ll=\"108288\"; __utmc=30149280; douban-fav-remind=1; ct=y; push_doumail_num=0; __utmz=30149280.1700924224.9.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmv=30149280.16703; push_noty_num=0; frodotk_db=\"f912931b88658561eb37b95299309023\"; ps=y; __utma=30149280.1262193242.1695303769.1701694458.1701699006.13; dbcl2=\"167037167:gpb9vz+25eM\"; ck=Qbp3; __utma=81379588.1626991796.1701702536.1701702536.1701702536.1; __utmc=81379588; __utmz=81379588.1701702536.1.1.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/misc/sorry; _pk_ref.100001.3ac3=%5B%22%22%2C%22%22%2C1701702536%2C%22https%3A%2F%2Fwww.douban.com%2Fmisc%2Fsorry%3Foriginal-url%3Dhttps%3A%2F%2Fbook.douban.com%2F%22%5D; _pk_id.100001.3ac3=100bef15360db890.1701702536.; _pk_ses.100001.3ac3=1; _vwo_uuid_v2=D861FF13AEF1F55EACED7545A6B63D3D4|9cc8c7e5fea2d936ec6102b68124c80b; __yadk_uid=HUNSEAjLbN6sAEFlcwhdBT9ayaw0WyzH; __utmt=1; __utmt_douban=1; __utmb=30149280.62.5.1701702423014; __utmb=81379588.10.10.1701702536"
)

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     "douban_book_list",
		WaitTime: 5 * time.Second,
		MaxDepth: 5,
		Cookie:   cookie,
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					URL:      "https://book.douban.com",
					Method:   "GET",
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": {ParseFunc: ParseTag},
			"书籍列表":  &collect.Rule{ParseFunc: ParseBookList},
			"书籍简介": &collect.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var public = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *collect.RuleContext) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))

	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoRe),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, public),
		"得分":  ExtraString(ctx.Body, scoreRe),
		"价格":  ExtraString(ctx.Body, priceRe),
		"简介":  ExtraString(ctx.Body, intoRe),
	}
	data := ctx.Output(book)

	result := collect.ParseResult{
		Items: []interface{}{data},
	}

	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {

	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

const BooklistRe = `<a.*?href="([^"]+)" title="([^"]+)"`

func ParseBookList(ctx *collect.RuleContext) (collect.ParseResult, error) {
	re := regexp.MustCompile(BooklistRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			URL:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &collect.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		result.Requests = append(result.Requests, req)
	}

	return result, nil
}

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *collect.RuleContext) (collect.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		if len(result.Requests) > 10 {
			break
		}
		result.Requests = append(
			result.Requests, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				URL:      fmt.Sprintf("%v%v", "https://book.douban.com", string(m[1])),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}
	return result, nil
}
