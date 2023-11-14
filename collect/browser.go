package collect

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/text/transform"
)

type BrowserFetch struct {
	Timeout time.Duration
}

func (fetch *BrowserFetch) Get(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: fetch.Timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed:%v", err)
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
