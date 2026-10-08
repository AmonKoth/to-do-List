package main

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Notes     string    `json:"notes"`
	Completed bool      `json:"completed"`
	Priority  int64     `json:"priority"`
	DueDate   time.Time `json:"due_date"`
}
