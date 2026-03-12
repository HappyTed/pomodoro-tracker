package deamon

import "context"

type (
	Server interface {
		Run(ctx context.Context) error
		Wait() error
	}

	Deamon interface {
		// Добавить задачу
		Add(name string, count int) error
		// Запустить таймер
		Run(ctx context.Context) error
		// Приостановить
		Pause() error
		// Завершить помидорку
		Stop() error
		// Сбросить все этапы
		Reset() error
		// Текущий статус по задаче
		Status() (Task, error)
		Shutdown() error // gracefull shutdown
	}

	Timer interface {
		Start(ctx context.Context) error
		Pause() error
		Reset() error
	}
)
