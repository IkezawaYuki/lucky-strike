package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"time"
)

var s *semaphore.Weighted = semaphore.NewWeighted(2)

func longProcess(ctx context.Context) {
	if err := s.Acquire(ctx, 1); err != nil {
		fmt.Println(err)
		return
	}
	defer s.Release(1)
	fmt.Println("Wait...")
	time.Sleep(1 * time.Second)
	fmt.Println("Done...")
}

func main() {
	ctx := context.Background()
	go longProcess(ctx)
	go longProcess(ctx)
	go longProcess(ctx)
	go longProcess(ctx)
	time.Sleep(5 * time.Second)
}
