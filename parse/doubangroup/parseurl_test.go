package doubangroup

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/donghc/crawler/collect"
)

func TestGetContent(t *testing.T) {
	cookie := "bid=qZ8gG_-P5M0; _pk_id.100001.8cb4=a4aa81c52ee51832.1691651193.; __yadk_uid=gfuYKDTSWELTL3H7yJ78VKbqNVJFAQeG; ll=\"108288\"; __utmz=30149280.1696663356.5.5.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; viewed=\"1007305\"; ap_v=0,6.0; __utmc=30149280; push_noty_num=0; push_doumail_num=0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1700120445%2C%22https%3A%2F%2Ftime.geekbang.com%2Fcolumn%2Farticle%2F612328%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.1471592734.1691651195.1700116287.1700120447.9; __utmt=1; __utmt_douban=1; loc-last-index-location-id=\"108288\"; _vwo_uuid_v2=DCBB38EF9DF817F80715096C3D1EE293B|d2f167ec2e012b3f6a02c950e5ce9b4d; dbcl2=\"167037167:j9y6vSfLpN8\"; ck=iZmW; __utmv=30149280.16703; __utmb=30149280.19.3.1700120488122"

	request := collect.Request{
		URL:       "https://www.douban.com/group/topic/297866960/?_i=0147635tC5jgSh",
		Cookie:    cookie,
		ParseFunc: ParseURL,
	}
	fmt.Println(request)
	content := GetContent([]byte("s"), "")

	fmt.Println(content.Items)
}

var (
	//go:embed data/1.txt
	s1 []byte
	//go:embed data/2.txt
	s2 []byte
)

func TestGetContentByStr(t *testing.T) {

	content := GetContent(s1, "tttt")

	fmt.Println(content.Items)

	content2 := GetContent(s2, "tttt")

	fmt.Println(content2.Items)
}
