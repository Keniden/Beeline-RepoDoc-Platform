package main

import "fmt"

type countryCodes map[string]int

type World struct {
	Name      string
	countries countryCodes
}

func NewWordl(name string) *World {
	return &World{
		Name:      name,
		countries: make(countryCodes),
	}
}

func (w *World) AddCountry(name string, code int) bool {
	w.countries[name] = code
	return true
}
  
func (w *World) GetCountryCode(name string) (code int, ok bool) {
	code, ok = w.countries[name] 
	return code, ok
}
// func main() {
// 	w := NewWordl("Earth")
// 	w.AddCountry("Russia", 52)
// 	fmt.Println(w.GetCountryCode("Russia"))

// }
