package deamon

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Task struct {
	name      string
	pomodoros int
	current   int
	is_active bool
}

var EMPTY_TASK = errors.New("empty task")

type TaskManager struct {
	// logger
	serv    Server
	mu      sync.RWMutex
	wg      sync.WaitGroup
	task    *Task
	timer   Timer
	signals chan struct{}
}

func NewDeamon(timer time.Duration) (Deamon, error) {
	s := make(chan struct{})
	t, err := NewTimer(timer, s)
	if err != nil {
		return nil, err
	}
	tm := &TaskManager{timer: t, signals: s}
	return tm, nil
}

// Добавить задачу
func (tm *TaskManager) Add(name string, count int) error {
	if tm.task != nil {
		return errors.New("Задача уже заведена, невозможно добавить новую")
	}
	t := &Task{
		name:      name,
		pomodoros: count,
		is_active: false,
		current:   0,
	}
	tm.task = t
	return nil
}

// Запустить таймер
func (tm *TaskManager) Run(ctx context.Context) error {
	if err := tm.timer.Start(ctx); err != nil {
		return err
	}
	tm.wg.Add(1)
	go func() {
		defer tm.wg.Done()
		select {
		case <-ctx.Done():
			return
		case <-tm.signals:
			tm.mu.Lock()
			defer tm.mu.Unlock()

			tm.task.current++
			tm.task.is_active = false

			return
		}
	}()

	tm.mu.Lock()
	tm.task.is_active = true
	tm.mu.Unlock()

	return nil
}

// Приостановить
func (tm *TaskManager) Pause() error {
	tm.mu.Lock()
	tm.task.is_active = false
	tm.mu.Unlock()

	tm.timer.Pause()

	return nil
}

// Завершить помидорку
func (tm *TaskManager) Stop() error {
	tm.mu.Lock()
	tm.task.is_active = false
	tm.task.current += 1
	tm.mu.Unlock()

	tm.timer.Reset()

	return nil
}

// Сбросить все этапы
func (tm *TaskManager) Reset() error {
	tm.mu.Lock()
	tm.task.is_active = false
	tm.task.current = 0
	tm.mu.Unlock()

	tm.timer.Reset()

	return nil
}

// Информация по задаче
func (tm *TaskManager) Status() (Task, error) {
	if tm.task == nil {
		return Task{}, EMPTY_TASK
	}
	return *tm.task, nil
}

func (tm *TaskManager) Shutdown() error {
	tm.wg.Wait()
	return nil
}
