package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func CreateTask(id int64, name string, notes string, priority int64, dueDate time.Time) Task {
	return Task{
		ID:        id,
		Name:      name,
		Notes:     notes,
		Completed: false,
		Priority:  priority,
		DueDate:   dueDate,
	}
}

func NewTask(tasks []Task, id int64, name string, notes string, priority int64, dueDate time.Time) []Task {
	t := CreateTask(id, name, notes, priority, dueDate)
	fmt.Println(t)
	tasks = append(tasks, t)
	return tasks
}

func OpenList(fileName string) *os.File {
	jsonFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return jsonFile
}

func ReadTasks(f *os.File, tasks *[]Task) error {

	err := json.NewDecoder(f).Decode(tasks)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func SaveFile(f *os.File, tasks []Task) error {
	_, err := f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seek: %w", err)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	if tasks == nil {
		tasks = []Task{}
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", " ")
	return encoder.Encode(tasks)
}

func UpdateTask(taskID int64, task Task, taskList []Task) error {
	if taskID < int64(len(taskList)) && taskID >= 0 {
		taskList[taskID] = task
		return nil
	}
	return fmt.Errorf("entry not updated")
}

func DeleteTasks(taskList []Task, index int) ([]Task, error) {

	if index < 0 || index >= len(taskList) {
		return taskList, fmt.Errorf("task not found")
	}
	taskList = append(taskList[:index], taskList[index+1:]...)
	return ReorderList(taskList), nil
}

func PrintList(taskList []Task) {
	if len(taskList) == 0 {
		fmt.Println("Tasks is empty")
		return
	}
	for i := range taskList {
		fmt.Println(taskList[i])
	}
}

func ReorderList(taskList []Task) []Task {
	for i := range taskList {
		taskList[i].ID = int64(i + 1)
	}
	return taskList
}

func CloseList(f *os.File) error {
	if f == nil {
		return nil
	}
	err := f.Close()
	if err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil

}
