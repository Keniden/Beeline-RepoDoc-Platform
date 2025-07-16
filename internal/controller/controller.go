package controller

import (
	"to-do-list/internal/storage"
	"to-do-list/internal/tasks"
)



type Storage interface {
	Create(task tasks.Task) error
}

type Controller struct{
	storage Storage
}

func NewController(storage Storage) *Controller {
    return &Controller{
        storage: storage,
    }
}


