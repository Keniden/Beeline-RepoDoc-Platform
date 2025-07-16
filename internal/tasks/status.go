package tasks

type Status uint8

const (
	ToDO Status = iota
	InProgress
	Done
)
