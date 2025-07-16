package tasks


import "time"

type Task struct {
	Name string
	Description string
	Deadline time.Time
	Status Status
}


