package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"example.com/client"
)

func main() {
	c := client.NewClient()
	wg := sync.WaitGroup{}
	ctx := context.Background()
	for i := 0; i < 400; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
			defer wg.Done()
			c.GetSleepy(ctx)
		}()
	}
	wg.Wait()
	fmt.Println("sleeping")
	time.Sleep(30 * time.Second)
	fmt.Println("finished")
}
