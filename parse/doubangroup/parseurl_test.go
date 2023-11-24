package doubangroup

import (
	_ "embed"
	"fmt"
	"testing"
)

var (
	//go:embed data/1.txt
	s1 []byte
	//go:embed data/2.txt
	s2 []byte
)

func TestGetContent(t *testing.T) {

	content := GetContent(s1, "tttt")

	fmt.Println(content.Items)

	content2 := GetContent(s2, "tttt")

	fmt.Println(content2.Items)
}
