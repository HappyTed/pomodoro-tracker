package deamon

import (
	"context"
	"errors"
	"time"
)

type Task struct {
	name      string
	pomodoros int
	current   int
	is_active bool
}

type TaskManagerUC struct {
	// logger
	// timer
	task    *Task
	timer   Timer
	signals chan SIGNALS
}

func NewTaskManager(timer time.Duration) (TaskManager, error) {
	s := make(chan SIGNALS)
	t, err := NewTimer(timer, s)
	if err != nil {
		return nil, err
	}
	tm := &TaskManagerUC{timer: t, signals: s}
	return tm, nil
}

// Добавить задачу
func (tm *TaskManagerUC) Add(name string, count int) error {
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
func (tm *TaskManagerUC) Run(ctx context.Context) error {
	err := tm.timer.Start(ctx)
	go func() {
		select {
		case sig := <-tm.signals:
			if sig == DONE {
				tm.task.current++
				return
			}
		case <-ctx.Done():
			return
		}
	}()
	return err
}

// Приостановить
func (tm *TaskManagerUC) Pause() error {
	tm.timer.Pause()
	return nil
}

// Завершить помидорку
func (tm *TaskManagerUC) Stop() error {
	tm.timer.Reset()
	tm.task.current += 1
	return nil
}

// Сбросить все этапы
func (tm *TaskManagerUC) Reset() error {
	tm.timer.Reset()
	tm.task.current = 0
	return nil
}

// Текущий статус
func (tm *TaskManagerUC) Status() error {
	return nil
}
