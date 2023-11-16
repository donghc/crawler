package doubangroup

import (
	"regexp"

	"github.com/donghc/crawler/collect"
)

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(cityListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(
			result.Requests, &collect.Request{
				URL:    u,
				Cookie: req.Cookie,
				ParseFunc: func(c []byte, request *collect.Request) collect.ParseResult {
					return GetContent(c, u)
				},
			})
	}
	return result
}

// #link-report > div
// /html/body/div[3]/div[1]/div/div[1]/div[2]/div[2]/div[1]/div
//const ContentRe = `<div class="rich-content topic-richtext">[\s\S]*?阳台[\s\S]*?</div>`

const ContentRe = `<div\s+class="topic-content">(?s:.)*?</div>`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)
	resultStr := re.FindString(string(contents))

	r2 := regexp.MustCompile("阳台")

	ok := r2.MatchString(resultStr)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	result := collect.ParseResult{
		Items: []interface{}{url},
	}

	return result
}
