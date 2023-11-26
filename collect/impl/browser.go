package impl

import (
	"bufio"
	"fmt"
	"github.com/donghc/crawler/collect"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/text/transform"

	"github.com/donghc/crawler/extensions"
	"github.com/donghc/crawler/proxy"
)

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.ProxyFunc
}

func (fetch *BrowserFetch) Get(request *collect.Request) ([]byte, error) {
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

	if request.Task.Cookie != "" {
		req.Header.Set("Cookie", request.Task.Cookie)
	}
	req.Header.Set("User-Agent", extensions.GenerateRandomUA())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := collect.DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}
