package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"

	"github.com/donghc/crawler/collect"
	"github.com/donghc/crawler/parse/doubangroup"
)

var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	douban()
}

func doubanGroup() {
	var workList []*collect.Request
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("<https://www.douban.com/group/szsh/discussion?start=%d>", i)
		workList = append(workList, &collect.Request{
			Ulr:       str,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	// collect.BrowserFetch{Timeout: 3 * time.Second}

}

func douban() {
	// url := "https://www.thepaper.cn/"
	url := "https://book.douban.com/subject/1007305/"
	f := collect.BrowserFetch{
		Timeout: 10 * time.Second,
	}

	body, err := f.Get(url)

	if err != nil {
		fmt.Printf("read content failed:%v \n", err)
		return
	}
	fmt.Println(string(body))
}

func test() {
	// 创建谷歌浏览器实例
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	// 创建一个执行器（allocator），使用选项列表
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	// 创建一个浏览器实例
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	defer browserCancel()

	ctx, cancel := context.WithTimeout(browserCtx, 15*time.Second)
	defer cancel()

	// 爬取页面，等待某一个元素出现，接着模拟鼠标点击，最后获取数据
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\\n%s", example)
	time.Sleep(5 * time.Second)

}
