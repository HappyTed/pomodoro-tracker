package deamon

import (
	"context"
	"fmt"
	"net"
	"os"
)

func handleCommand(c net.Conn) error {
	return nil
}

// Unix Socker Server
type Server struct {
	lst net.Listener
}

func New(socketPath string) (*Server, error) {
	var s = &Server{}

	_, err := os.Stat(socketPath)
	if err != nil {
		// сокет существует
	}

	if err == nil {
		// если файл сокета уже существует, нужно его удалить для дальнейшей работы программы
		fmt.Println("Deleting existing", socketPath)
		err := os.Remove(socketPath)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("listen error:", err)
		return nil, err
	}
	s.lst = l

	return s, nil
}

func (s *Server) Run(ctx context.Context) {
	for {
		nc, err := s.lst.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			return
		}
		go handleCommand(nc)
	}
}
