package doubangroup

import (
	_ "embed"
	"fmt"
	"github.com/donghc/crawler/collect"
	"testing"
)

var (
	//go:embed data/1.txt
	s1 []byte
	//go:embed data/2.txt
	s2 []byte
)

func TestGetContent(t *testing.T) {

	content, _ := GetSunRoom(&collect.RuleContext{
		Body: s1,
		Req:  nil,
	})

	fmt.Println(content.Items)

	content2, _ := GetSunRoom(&collect.RuleContext{
		Body: s2,
		Req:  nil,
	})

	fmt.Println(content2.Items)
}
