package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/context"

	"github.com/donghc/crawler/collect"
)

var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	main2()
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

func main2() {
	ctx := context.Background()
	before := time.Now()
	preCtx, _ := context.WithTimeout(ctx, 500*time.Millisecond)
	go func() {
		childCtx, _ := context.WithTimeout(preCtx, 300*time.Millisecond)
		select {
		case <-childCtx.Done():
			after := time.Now()
			fmt.Println("child during:", after.Sub(before).Milliseconds())
		}
	}()
	select {
	case <-preCtx.Done():
		after := time.Now()
		fmt.Println("pre during:", after.Sub(before).Milliseconds())
	}
}
