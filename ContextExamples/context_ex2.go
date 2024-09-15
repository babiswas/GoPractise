package ContextExamples

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func CancelContext() {
	var wg sync.WaitGroup

	fmt.Println("Displaying go context by cancelling it.")
	ctx, cancel := context.WithCancel(context.Background())

	message := make(chan string, 3)
	confirm := make(chan string)
	cancel_context := make(chan string)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(ctx context.Context, value string, ch chan<- string, wg *sync.WaitGroup) error {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					time.Sleep(1 * time.Millisecond)
				}
			}
			ch <- value
			return nil
		}(ctx, "message:"+strconv.Itoa(i), message, &wg)
	}

	cancel_ctx := func(ch chan<- string) {
		time.Sleep(3 * time.Millisecond)
		cancel()
		ch <- "Cancelled context after 3 millisecond."
	}

	message_collector := func(confirm chan<- string, result <-chan string) {
		for value := range result {
			fmt.Println(value)
		}
		confirm <- "All messages processed."
	}

	go message_collector(confirm, message)
	go cancel_ctx(cancel_context)

	wg.Wait()
	close(message)
	fmt.Println(<-cancel_context)
	fmt.Println(<-confirm)
}
