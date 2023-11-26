package impl

import (
	"bufio"
	"fmt"
	"github.com/donghc/crawler/collect"
	"io/ioutil"
	"net/http"

	"golang.org/x/text/transform"
)

type SimpleFetch struct {
}

func (fetch *SimpleFetch) Get(req *collect.Request) ([]byte, error) {
	resp, err := http.Get(req.URL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%d \n", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := collect.DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}
