package deamon

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskManagerEmptyTask(t *testing.T) {
	tm, err := NewDeamon(TIMER_LENGTH)
	defer tm.Shutdown()
	if err != nil {
		t.Fatal("Не удалось инициализировать таск-менеджер из-за ошибки:", err)
		return
	}

	_, err = tm.Status()

	if err != EMPTY_TASK {
		t.Error("До создания задачи, task manager должен хранить nil указатель")
	}
}

func TestTaskManagerAddTask(t *testing.T) {
	tm, err := NewDeamon(TIMER_LENGTH)
	defer tm.Shutdown()
	if err != nil {
		t.Fatal("Не удалось инициализировать таск-менеджер из-за ошибки:", err)
		return
	}

	name := "новая тестовая задача"
	count := 4
	tm.Add(name, count)

	actual, err := tm.Status()
	if err != nil {
		t.Error("Ошибка получения информации о задаче:", err)
	}

	expected := Task{
		name:      name,
		pomodoros: count,
		current:   0,
		is_active: false,
	}

	assert.Equal(t, expected, actual)
}

func TestTaskManagerRunPomodoro(t *testing.T) {
	tm, err := NewDeamon(TIMER_LENGTH)
	defer tm.Shutdown()
	if err != nil {
		t.Fatal("Не удалось инициализировать таск-менеджер из-за ошибки:", err)
		return
	}

	name := "новая тестовая задача"
	count := 4
	tm.Add(name, count)

	actual, err := tm.Status()
	if err != nil {
		t.Error("Ошибка получения информации о задаче:", err)
	}

	expected := Task{
		name:      name,
		pomodoros: count,
		current:   0,
		is_active: false,
	}

	assert.Equal(t, expected, actual) // если ошибка, дальше нет смысла идти

	ctx, cancel := context.WithCancel(context.Background())
	tm.Run(ctx)

	// Проверить, что задача стала активной
	time.Sleep(1 * time.Second)

	actual, err = tm.Status()
	if err != nil {
		t.Error("Ошибка получения информации о задаче:", err)
	}

	expected = Task{
		name:      name,
		pomodoros: count,
		current:   0,
		is_active: true,
	}

	assert.Equal(t, expected, actual)

	// Проверить, что задача стала неактивной и счётчик увеличился
	time.Sleep(TIMER_LENGTH)

	actual, err = tm.Status()
	if err != nil {
		t.Error("Ошибка получения информации о задаче:", err)
	}

	expected = Task{
		name:      name,
		pomodoros: count,
		current:   1,
		is_active: false,
	}

	assert.Equal(t, expected, actual)

	cancel()
}
