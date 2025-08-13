// package main

// import (
// 	"context"
// 	"fmt"
// )

// func generate(ctx context.Context, out chan<- int, n int) {
// 	defer close(out)

// 	for i := 2; i < n+1; i++ {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case out <- i:
// 		}
// 	}
// }

// func filterEven(ctx context.Context, in <-chan int, out chan<- int) {
// 	defer close(out)
// 	for v := range in {
// 		if v%2 == 0 {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case out <- v:
// 			}
// 		}
// 	}
// }

// func square(ctx context.Context, in <-chan int, out chan<- int) {
// 	defer close(out)
// 	for v := range in {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case out <- v * v:
// 		}
// 	}
// }

// func main() {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	var ch1 = make(chan int, 4)
// 	var ch2 = make(chan int, 4)
// 	var ch3 = make(chan int, 4)
// 	go generate(ctx, ch1, 50)
// 	go filterEven(ctx, ch1, ch2)
// 	go square(ctx, ch2, ch3)
// 	for v := range ch3 {
// 		if v > 50 {
// 			cancel()
// 		}
// 		fmt.Println(v)

// 	}

// }
