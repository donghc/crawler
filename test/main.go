package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()
	script := `
    var n = 100
    console.log("hello-" + n)
    n = n-10;
    n-10
  `
	value, _ := vm.Run(script)
	fmt.Println("value:", value.String())
}
