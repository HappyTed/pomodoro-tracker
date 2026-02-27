package entities

import "time"

type Task struct {
	Name     string
	State    uint
	IsActive bool
	Target   int
	Current  int
	Timer    time.Duration
}
