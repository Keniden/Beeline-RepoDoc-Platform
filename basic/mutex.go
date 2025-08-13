package main

import (
	"fmt"
	"sync"
)

type LinkStorage struct {
	sync.Mutex
	storage map[string]string
}

func (s *LinkStorage) SetValue(key, value string) {
	s.Lock()
	defer s.Unlock()
	s.storage[key] = value
}

func NewLinkStorage() *LinkStorage{
	return &LinkStorage{
		storage: make(map[string]string),
	}
}
// func main(){
// 	storage := NewLinkStorage()
// 	storage.SetValue("Россия", "Абоба")
// 	fmt.Println(storage.storage)
// }

func init(){
	
}

func worlds(i int) int{
	return i + i
}
