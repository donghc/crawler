package doubangroup

import (
	"regexp"

	"github.com/donghc/crawler/collect"
)

const cityListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte) collect.ParseResult {
	re := regexp.MustCompile(cityListRe)

	matches := re.FindAllSubmatch(contents, -1)

	parseResult := collect.ParseResult{}

	for _, v := range matches {
		url := string(v[1])

		parseResult.Requests = append(parseResult.Requests, &collect.Request{
			Ulr: url,
			ParseFunc: func(c []byte) collect.ParseResult {
				return GetContent(c, url)
			},
		})
	}

	return parseResult
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)

	ok := re.Match(contents)
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
