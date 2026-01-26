package usecase

type IPomodoro interface {
	Run() error    // Запуск
	Cancel() error // Остановить (аналогично stop)
	Pause() error  // Поставить на пазу
	Reset() error  // Сбросить таймер
	Current() (int, error)
}
