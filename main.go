package main

import (
	"GoPractise/ContextExamples"
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("==========================")
	fmt.Println("Displaying timeout context")
	ContextExamples.TimeoutContext()
	fmt.Println("==========================")
	fmt.Println("Displaying go context cancelling it after sometime.")
	ContextExamples.CancelContext()
	fmt.Println("===========================")
	fmt.Println("Another timeout example:")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	count, err := ContextExamples.Webapi(ctx)
	if err != nil {
		fmt.Println("Long running task exited with error:", err)
		return
	}
	fmt.Println("Count is:", count)
	fmt.Println("============================")
	fmt.Println("Webserver:")
	ContextExamples.Webserver()
}
