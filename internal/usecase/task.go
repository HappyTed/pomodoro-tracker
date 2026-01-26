package usecase

import (
	"github.com/google/uuid"
)

type Counter struct {
	Max     uint8
	Current uint8
}

type Task struct {
	Id uuid.UUID
	// Type        string
	Name        string
	Description string
	State       bool // false - inactive
	Counter

	tomato IPomodoro
}

type TaskOption func(*Task) error

func WithTaskName(name string) TaskOption {
	return func(t *Task) error {
		t.Name = name
		return nil
	}
}

func WithTaskDescription(description string) TaskOption {
	return func(t *Task) error {
		t.Description = description
		return nil
	}
}

func WithTomatoUC(tomato IPomodoro) TaskOption {
	return func(t *Task) error {
		t.tomato = tomato
		return nil
	}
}

func NewTask(options ...TaskOption) (*Task, error) {
	task := &Task{
		Id: uuid.New(),
	}

	for _, opt := range options {
		err := opt(task)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}
