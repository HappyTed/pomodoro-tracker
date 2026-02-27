package deamon

import (
	"context"
	"fmt"
	"time"
)

type TimerUC struct {
	target    time.Duration // длительность таймера
	remaining time.Duration // оставшееся время
	start_at  time.Time     // время запуска
	breakCh   chan struct{} // внутренний сигнал на паузу или сброс
	signal    chan struct{} // сигнал во вне об успешном завершении
}

func NewTimer(target time.Duration, signal chan struct{}) (Timer, error) {
	t := &TimerUC{
		target:    target,
		remaining: target,
		breakCh:   make(chan struct{}),
		signal:    signal,
		start_at:  time.Time{},
	}
	return t, nil
}

func (t *TimerUC) Start(ctx context.Context) error {
	fmt.Println("Оставшееся время:", t.remaining)
	go func() {
		t.start_at = time.Now()
		timeoutCh := time.After(t.remaining)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.breakCh:
				return
			case <-timeoutCh:
				t.signal <- struct{}{}
				t.start_at = time.Time{}
				t.remaining = t.target
				return
			}
		}
	}()

	return nil
}

func (t *TimerUC) Pause() error {
	t.breakCh <- struct{}{}
	stop := time.Now()
	workTime := stop.Sub(t.start_at)
	t.remaining = (t.target - workTime).Round(time.Second)
	return nil
}

func (t *TimerUC) Reset() error {
	t.breakCh <- struct{}{}
	t.start_at = time.Time{}
	t.remaining = t.target
	return nil
}

func (t *TimerUC) Status() time.Duration {
	return t.remaining
}
