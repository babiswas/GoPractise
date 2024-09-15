package ContextExamples

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func TimeoutContext() {
	fmt.Println("Displaying the timeout context.")
	var wg sync.WaitGroup

	confirm := make(chan string)
	message := make(chan string, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(ctx context.Context, value string, wg *sync.WaitGroup, ch chan<- string) error {
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
		}(ctx, "message"+strconv.Itoa(i), &wg, message)
	}

	message_collector := func(result <-chan string, confirm chan<- string) {
		for value := range result {
			fmt.Println(value)
		}
		confirm <- "All messages has been processed."
	}

	go message_collector(message, confirm)

	wg.Wait()
	close(message)
	fmt.Println(<-confirm)

}
