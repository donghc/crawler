package doubangroup

import (
	"fmt"
	"testing"

	"github.com/donghc/crawler/collect"
)

func TestGetContent(t *testing.T) {
	request := collect.Request{
		URL:       "",
		Cookie:    "",
		ParseFunc: ParseURL,
	}
	fmt.Println(request)
	content := GetContent([]byte("s"), "")

	fmt.Println(content.Items)
}
