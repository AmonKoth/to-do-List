package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	var tasks []Task
	var file *os.File

	reader := bufio.NewReader(os.Stdin)
	running := true

	FunctionText()

	for running {
		operation := GetInput(reader)
		operation = strings.ToLower(operation)
		if operation == "exit" || operation == "e" {
			HandleSave(file, tasks)
			fmt.Println("exiting")
			running = false
		}
		switch operation {
		case "open", "o":
			if file != nil {
				HandleSave(file, tasks)
				err := file.Close()
				WriteupError(err)
			}
			fmt.Println(":opening file")
			file = HandleOpen(reader)
			if file != nil {
				fmt.Println("File open")
				fmt.Println(":reading file")
				_, err := file.Seek(0, 0)
				WriteupError(err)
				tasks = nil
				err = ReadTasks(file, &tasks)
				WriteupError(err)
			} else {
				tasks = nil
			}

		case "print", "p":
			PrintList(tasks)
		case "create", "c", "create task":
			tasks = HandleCreate(tasks, reader)
		case "update task", "u", "update":
			HandleUpdate(tasks, reader)
		case "delete", "d", "delete task":
			listAfterDelete, err := HandleDelete(tasks, reader)
			if err != nil {
				WriteupError(err)
				continue
			}
			tasks = listAfterDelete
		}
	}
	if file != nil {
		WriteupError(CloseList(file))
	}
}

func HandleOpen(reader *bufio.Reader) *os.File {
	fmt.Println("Give the file name")
	input := GetInput(reader)
	if input == "" {
		return nil
	}
	fmt.Println("read string is: ", input)
	return OpenList(input)
}

func HandleCreate(tasks []Task, reader *bufio.Reader) []Task {

	id := int64(len(tasks) + 1)
	fmt.Println("Give the name of the task")
	name := GetInput(reader)
	fmt.Println("Give the notes for the task")
	notes := GetInput(reader)
	fmt.Println("Give priority for the task")
	input := GetInput(reader)
	priority, err := strconv.ParseInt(input, 10, 64)
	WriteupError(err)
	if err != nil {
		return tasks
	}
	fmt.Println("Give the due date in DD-MM-YYYY")
	dateInput := GetInput(reader)
	dueDate, err := time.Parse("02-01-2006", dateInput)
	WriteupError(err)
	if err != nil {
		return tasks
	}
	tasks = NewTask(tasks, id, name, notes, priority, dueDate)
	FunctionText()
	return tasks
}

func HandleUpdate(taskList []Task, reader *bufio.Reader) {
	fmt.Println("Enter the id of the task to be edited ")
	taskID := GetInput(reader)
	id, err := strconv.ParseInt(taskID, 10, 64)
	WriteupError(err)
	if err != nil {
		return
	}
	for i := range taskList {
		if int64(i+1) != id {
			continue
		}
		foundTask := taskList[i]

		updating := true
		for updating {
			fmt.Println("Select what you want to update 1: name 2: notes 3: priority 4: due date 5: task complete? e: end")
			operation := GetInput(reader)

			switch operation {
			case "1", "name":
				fmt.Println("Give the new name of the task ")
				foundTask.Name = GetInput(reader)
			case "2", "notes":
				fmt.Println("Give the updated notes for the task")
				foundTask.Notes = GetInput(reader)
			case "3", "p", "priority":
				fmt.Println("Give the updated priority")
				priorityInput := GetInput(reader)
				priority, err := strconv.ParseInt(priorityInput, 10, 64)
				WriteupError(err)
				if err == nil {
					foundTask.Priority = priority
				}
			case "4", "d", "date":
				fmt.Println("Give the new due date of the task")
				dateInput := GetInput(reader)
				dueDate, err := time.Parse("02-01-2006", dateInput)
				WriteupError(err)
				if err == nil {
					foundTask.DueDate = dueDate
				}
			case "5", "c", "complete", "done":
				fmt.Println("Is the task done?")
				isDone := GetInput(reader)
				isDone = strings.ToLower(isDone)
				switch isDone {
				case "yes", "y", "t", "true":
					foundTask.Completed = true
				default:
					foundTask.Completed = false
				}
			case "e", "end":
				updating = false
			}
		}
		err = UpdateTask(int64(i), foundTask, taskList)
		WriteupError(err)
		return
	}

	fmt.Println("Task with the given ID cannot be found")
	FunctionText()
	return
}

func GetInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		return "exit"
	}
	WriteupError(err)
	return strings.TrimSpace(input)
}
func WriteupError(err error) {
	if err != nil {
		fmt.Println("Error Occurred: ", err)
	}
}

func HandleSave(file *os.File, tasks []Task) {
	if file == nil {
		return
	}
	err := SaveFile(file, tasks)
	WriteupError(err)
}
func HandleDelete(taskList []Task, reader *bufio.Reader) ([]Task, error) {
	if len(taskList) == 0 {
		return taskList, fmt.Errorf("task list is empty")
	}
	fmt.Println("Give the id of the item to be deleted")
	deleteInput := GetInput(reader)
	deleteId, err := strconv.Atoi(deleteInput)
	if err != nil {
		return taskList, err
	}
	fmt.Println("Deletion successful")
	return DeleteTasks(taskList, deleteId-1)
}

func FunctionText() {
	fmt.Println(`Usage:
Open: to open a file
Create: to create a new task
Print: to print the read file
Update: to update an existing task
Delete: to delete an existing task`)
}
