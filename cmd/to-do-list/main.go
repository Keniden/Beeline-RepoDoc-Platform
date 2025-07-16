package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"to-do-list/internal/controller"
	"to-do-list/internal/storage"
	"to-do-list/internal/tasks"

	flag "github.com/spf13/pflag"
)

var action string
var taskName string
var taskDescription string
var taskDeadline string

func init() {
	flag.StringVarP(&action, "action", "a", "", "тип задачи")

	flag.StringVarP(&taskName, "name", "n", "", "Наименование задачи")
	flag.StringVarP(&taskDescription, "description", "d", "", "Описание задачи")
	flag.StringVarP(&taskDeadline, "dd", "", "", "Дедлайн задачи Формат: 2006-01-02")

}

func main() {
	flag.Parse()
	log.Println(taskName, taskDescription, taskDeadline)
	// Создание файла
	const path = "db/storage.json"
	file, err := os.Open(path)

	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			log.Println(err)
		}
		file, err = os.Create(path)
		if err != nil {
			log.Println(err)
		}
	}

	controller.NewController(storage.NewStorage(file))

}
