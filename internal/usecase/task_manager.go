package usecase

import (
	"errors"
	"log"
)

type TaskManager struct {
	tasks  map[string]*Task
	target *Task
}

func NewTaskManager() (*TaskManager, error) {
	return &TaskManager{make(map[string]*Task), nil}, nil
}

func (tm *TaskManager) Add(data ...TaskOption) (string, error) {

	t, err := NewTask(data...)
	log.Printf("Try added new task with data: %+v\n", t)

	if err != nil {
		log.Printf("Error: %w", err)
		return "", err
	}

	uid := t.Id.String()
	tm.tasks[uid] = t

	log.Printf("Task uid is: %s\n", uid)

	tm.target = t

	return uid, err
}

func (tm *TaskManager) Remove(id string) error {
	if _, ok := tm.tasks[id]; !ok {
		return errors.New("invalid task id")
	}

	delete(tm.tasks, id)
	return nil
}

func (tm *TaskManager) Get(id string) (*Task, error) {
	if t, ok := tm.tasks[id]; ok {
		return t, nil
	} else {
		return nil, WRONG_ID
	}
}

func (tm *TaskManager) List() ([]*Task, error) {
	log.Println("Try get tasks list...")
	var result []*Task
	for _, val := range tm.tasks {
		result = append(result, val)
	}

	log.Printf("Tasks list is: %+v", result)

	return result, nil
}

func (tm *TaskManager) Target() (*Task, error) {
	return tm.target, nil
}
