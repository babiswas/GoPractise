package ContextExamples

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type output struct {
	count         int
	error_message error
}

func dbtask1(ctx context.Context, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-time.After(7 * time.Millisecond):
		return 30, nil
	}
}

func dbtask2(ctx context.Context, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		fmt.Println("DB2 access failure.")
		return 0, ctx.Err()
	case <-time.After(5 * time.Millisecond):
		return 20, nil
	}
}

func Webapi(ctx context.Context) (int, error) {
	wg := sync.WaitGroup{}
	output1 := make(chan output, 1)
	output2 := make(chan output, 1)

	wg.Add(1)
	go func() {
		count, err := dbtask1(ctx, &wg)
		o := output{
			count:         count,
			error_message: err,
		}
		output1 <- o
	}()

	wg.Add(1)
	go func() {
		count, err := dbtask2(ctx, &wg)
		o := output{
			count:         count,
			error_message: err,
		}
		output2 <- o
	}()
	wg.Wait()
	o1 := <-output1
	if o1.error_message != nil {
		return 0, o1.error_message
	}
	o2 := <-output2
	if o2.error_message != nil {
		return 0, o2.error_message
	}
	result := o1.count + o2.count
	return result, nil
}
