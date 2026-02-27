package main

import (
	"context"
	"sync"

	"pomodoro.tracker/internal/deamon"
)

var (
	mu *sync.Mutex
)

func main() {
	serv, err := deamon.New("/tmp/my_socket", 1024, 1)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serv.Run(ctx)
}
