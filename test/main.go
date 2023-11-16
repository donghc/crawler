package main

import (
	"fmt"
	"regexp"
)

func main() {
	html := `
        <div class="topic-content">
          这套房子有
        </div>
        
        <div>
          其它内容阳台
        </div>
    `

	r := regexp.MustCompile(`<div\s+class="topic-content">(?s:.)*?</div>`)
	result := r.FindString(html)

	r2 := regexp.MustCompile("阳台")
	if r2.MatchString(result) {
		fmt.Println("包含阳台")
	} else {
		fmt.Println("不包含阳台")
	}
}
