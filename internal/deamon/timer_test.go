package deamon

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	TIMER_LENGTH = 5 * time.Second
)

type TestData struct {
	Name  string
	Timer Timer
}

func TestDefaultFinish(t *testing.T) {
	sig := make(chan struct{})
	timer, err := NewTimer(TIMER_LENGTH, sig)
	if err != nil {
		t.Fatal("Не удалось инициализировать таймер из-за ошибки:", err)
		return
	}

	startAt := time.Now()
	fmt.Println("Запуск таймера:", TIMER_LENGTH)
	fmt.Println("Время запуска:", startAt)
	timer.Start(context.Background())

	fmt.Println("Ждём завершения таймера...")
	timeout := time.After(TIMER_LENGTH + 1*time.Second)

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был завершен за установленное время:", TIMER_LENGTH)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		stopAt := time.Now()
		fmt.Println("Успешное завершение таймера")
		fmt.Println("Время завершения", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
	}

	startAt = time.Now()
	fmt.Println("\nПовторный запуск таймера:", TIMER_LENGTH)
	fmt.Println("Время запуска:", startAt)
	timer.Start(context.Background())

	fmt.Println("Ждём завершения таймера...")
	timeout = time.After(TIMER_LENGTH + 1*time.Second)

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был завершен за установленное время:", TIMER_LENGTH)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		stopAt := time.Now()
		fmt.Println("Успешное завершение таймера")
		fmt.Println("Время завершения", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	}
}

func TestWithPause(t *testing.T) {
	sig := make(chan struct{})
	timer, err := NewTimer(TIMER_LENGTH, sig)
	if err != nil {
		t.Fatal("Не удалось инициализировать таймер из-за ошибки:", err)
		return
	}

	startAt := time.Now()
	fmt.Println("Запуск таймера:", TIMER_LENGTH)
	fmt.Println("Время запуска:", startAt)
	timer.Start(context.Background())

	timeout := time.After(TIMER_LENGTH)

	time.Sleep(1 * time.Second)
	fmt.Println("приостановка таймера через 1 секунду")
	timer.Pause()

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был приостановлен при нажатии на кнопку паузы :", TIMER_LENGTH)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		fmt.Println("ОК: таймер поставлен на паузу, время работы: 1 секунда")
	}

	startAt = time.Now()
	leftTime := TIMER_LENGTH - 1*time.Second
	timeout = time.After(TIMER_LENGTH)
	fmt.Println("Продолжаем работу таймера. Ожидаемое оставшееся время:", leftTime)
	timer.Start(context.Background())

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был завершен за оставшееся время:", leftTime)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		stopAt := time.Now()
		realWorkTime := (stopAt.Sub(startAt)).Round(time.Second)

		if realWorkTime != leftTime {
			t.Error("Время работы таймера привышет ожидаемое:", realWorkTime, "против", leftTime)
		}
		fmt.Println("Успешное завершение таймера")
		fmt.Println("Время завершения", stopAt)
		fmt.Println("Время работы:", realWorkTime)
		return
	}
}

func TestReset(t *testing.T) {
	sig := make(chan struct{})
	timer, err := NewTimer(TIMER_LENGTH, sig)
	if err != nil {
		t.Fatal("Не удалось инициализировать таймер из-за ошибки:", err)
		return
	}

	startAt := time.Now()
	fmt.Println("Запуск таймера:", TIMER_LENGTH)
	fmt.Println("Время запуска:", startAt)
	timer.Start(context.Background())

	timeout := time.After(TIMER_LENGTH)

	time.Sleep(1 * time.Second)
	fmt.Println("приостановка таймера через 1 секунду")
	timer.Pause()

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был приостановлен при нажатии на кнопку паузы :", TIMER_LENGTH)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		fmt.Println("ОК: таймер поставлен на паузу, время работы: 1 секунда")
	}

	startAt = time.Now()
	timeout = time.After(TIMER_LENGTH + 1*time.Second)
	fmt.Println("Сбрасываем таймер и снова запускаем его. Ожидаемое оставшееся время:", TIMER_LENGTH)
	timer.Reset()
	timer.Start(context.Background())

	select {
	case <-timeout:
		stopAt := time.Now()
		t.Error("Таймер не был завершен за оставшееся время:", TIMER_LENGTH)
		fmt.Println("Текущее время:", stopAt)
		fmt.Println("Время работы:", stopAt.Sub(startAt))
		return
	case <-sig:
		stopAt := time.Now()
		realWorkTime := (stopAt.Sub(startAt)).Round(time.Second)

		if realWorkTime != TIMER_LENGTH {
			t.Error("Время работы таймера привышет ожидаемое:", realWorkTime, "против", TIMER_LENGTH)
		}
		fmt.Println("Успешное завершение таймера")
		fmt.Println("Время завершения", stopAt)
		fmt.Println("Время работы:", realWorkTime)
		return
	}
}
