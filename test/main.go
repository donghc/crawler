package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	limit := rate.NewLimiter(rate.Every(500*time.Millisecond), 1)
	for {
		if err := limit.Wait(context.Background()); err == nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		}
	}
}
