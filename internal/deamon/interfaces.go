package deamon

import "context"

type (
	TaskManager interface {
		// Добавить задачу
		Add(name string, count int) error
		// Запустить таймер
		Run() error
		// Приостановить
		Pause() error
		// Завершить помидорку
		Stop() error
		// Сбросить все этапы
		Reset() error
		// Текущий статус
		Status() error
	}

	Timer interface {
		Start(ctx context.Context) error
		Pause() error
		Reset() error
	}
)
