package main

import (
	"fmt"
	"sync"
)


func main() {
	var mx sync.Mutex
	mx.Lock()
	defer func() {
		if r := recover(); r != nil{
			fmt.Println("panica")
		}
	mx.Unlock()
	}()
	panic("Pizda")
}

