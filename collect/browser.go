package collect

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/text/transform"

	"github.com/donghc/crawler/proxy"
)

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
}

func (fetch *BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: fetch.Timeout,
	}

	if fetch.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = fetch.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
	}

	if request.Cookie != "" {
		req.Header.Set("Cookie", request.Cookie)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
