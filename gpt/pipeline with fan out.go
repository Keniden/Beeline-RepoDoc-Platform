package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func generate(ctx context.Context, n int, out chan<- int) {
	defer close(out)
	for i := 1; i <= n; i++ {
		select {
		case <-ctx.Done():
			return
		case out <- i:
		}
	}

}

func process(ctx context.Context, in <-chan int, out chan<- int) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-in:
			if !ok {
				return
			}
			time.Sleep(100 * time.Millisecond)
			res := v * v
			select {
			case <-ctx.Done():
				return
			case out <- res:
			}
		}
	}
}

func merge(ctx context.Context, chs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	forward := func(ch <-chan int) {
		defer wg.Done()
		for v := range ch {
			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		}
	}

	for _, ch := range chs {
		go forward(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	M := 10
	in := make(chan int)
	outs := make([]chan int, M)
	for i := range outs {
		outs[i] = make(chan int, 4)
		go process(ctx, in, outs[i])
	}

	go generate(ctx, 100, in) // генератор запускаем после воркеров

	chs := make([]<-chan int, len(outs))
	for i, ch := range outs {
		chs[i] = ch
	}

	merged := merge(ctx, chs...) // без <-

	for v := range merged {
		fmt.Println(v)
	}
}
