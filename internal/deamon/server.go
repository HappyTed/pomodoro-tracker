package deamon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"pomodoro.tracker/internal/models/api"
)

// Unix Socker Server
type UnixSocketServer struct {
	// logger
	listner        net.Listener
	maxBufSize     int64 // kb, по дефолту можно сделать 128
	maxConnections int
	connections    int
	mu             sync.RWMutex
	socketPath     string
}

func New(socketPath string, bufSize int64, maxConnections int) (Server, error) {
	var s = &UnixSocketServer{
		maxBufSize:     bufSize,
		maxConnections: maxConnections,
		socketPath:     socketPath,
	}

	return s, nil
}

func (s *UnixSocketServer) Run(ctx context.Context) error {
	fmt.Println("Try run unix socket domain server")

	_, err := os.Stat(s.socketPath)
	if err == nil {
		// если файл сокета уже существует, нужно его удалить для дальнейшей работы программы
		fmt.Println("Deleting existing", s.socketPath)
		err := os.Remove(s.socketPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	l, err := net.Listen("unix", s.socketPath)
	if err != nil {
		fmt.Println("listen error:", err)
		return err
	}
	s.listner = l

	fmt.Printf(
		"Server run! Waiting for connection: %s\n",
		s.listner.Addr().String(),
	)

	handelConnections(ctx)

	return nil

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("Server stopped by context")
	// 		s.listner.Close()
	// 		return nil
	// 	default:
	// 		if s.connections < s.maxConnections {
	// 			nc, err := s.listner.Accept() // Это ожидающая операция
	// 			// и к этому момоенту мог сработать контекст, получается гонка
	// 			if err != nil {
	// 				fmt.Println("Accept error:", err)
	// 				continue
	// 			}
	// 			fmt.Println("New connection!")
	// 			s.mu.Lock()
	// 			s.connections++
	// 			s.mu.Unlock()
	// 			go s.handleCommand(ctx, nc)
	// 		}

	// 	}
	// }
}

func (s *UnixSocketServer) Stop() {
	if s.listner != nil {
		s.listner.Close()
	}
}

type HandlerFunc func(c net.Conn) error

func middleware(handler HandlerFunc) HandlerFunc {
	// Ограничение по max connections
	// логирование
	return func(c net.Conn) error {
		err := handler(c)
		return err
	}
}

func handlerFactory(ctx context.Context, c net.Conn, buffSize int64) {
	for {
		buf := make([]byte, buffSize)
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Close Connection by client.")
				return
			}
			fmt.Println("Read Error:", err)
			return
		}

		data := buf[:n]
		fmt.Print("Server got:", string(data))

		req := &api.Request{}
		resp := &api.Response{
			Status: api.OK,
		}
		err = json.Unmarshal(data, req)
		if err != nil {
			fmt.Println("Marshal read data Error:", err)
			resp = &api.Response{
				Status:  api.ERROR,
				Message: err.Error(),
			}
		}

		// TODO: где-то здесь вызывать нужны обработчик
		if cmd, ok := api.Commands[req.Cmd]; ok {
			switch cmd {
			case api.ADD:
				AddTaskHandleFunc()
			case api.START:
				StartHandleFunc()
			case api.STOP:
				StopHandleFunc()
			case api.PAUSE:
				PauseHandleFunc()
			case api.RESET:
				ResetHandleFunc()
			case api.STATUS:
				StatusHandleFunc()
			}
		} else {
			resp = &api.Response{
				Status:  api.ERROR,
				Message: fmt.Sprintf("Unknow Command: %s\n", req.Cmd),
			}
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			_, err = c.Write([]byte(err.Error()))
			fmt.Println("Failed to prepare json response: marshal error")
			return
		}

		_, err = c.Write(respData)
	}
}

// работает с подключением до обрыва соединения
func handelConnections(lst net.Listener, ch chan net.Listener) {

	for {

		conn, err := lst.Accept()
		if err != nil {
			// TODO
		}
		ch <- conn

		select {
		case <-ctx.Done():
			fmt.Println("Server stopped by context")
			return
		default: // не блокирующая операция
			conn, err := s.listner.Accept()

			buf := make([]byte, s.maxBufSize)
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Close Connection by client.")
					return
				}
				fmt.Println("Read Error:", err)
				return
			}

			data := buf[:n]
			fmt.Print("Server got:", string(data))

			req := &api.Request{}
			resp := &api.Response{
				Status: api.OK,
			}
			err = json.Unmarshal(data, req)
			if err != nil {
				fmt.Println("Marshal read data Error:", err)
				resp = &api.Response{
					Status:  api.ERROR,
					Message: err.Error(),
				}
			}

			// TODO: где-то здесь вызывать нужны обработчик
			if cmd, ok := api.Commands[req.Cmd]; ok {
				switch cmd {
				case api.ADD:
					AddTaskHandleFunc()
				case api.START:
					StartHandleFunc()
				case api.STOP:
					StopHandleFunc()
				case api.PAUSE:
					PauseHandleFunc()
				case api.RESET:
					ResetHandleFunc()
				case api.STATUS:
					StatusHandleFunc()
				}
			} else {
				resp = &api.Response{
					Status:  api.ERROR,
					Message: fmt.Sprintf("Unknow Command: %s\n", req.Cmd),
				}
			}

			respData, err := json.Marshal(resp)
			if err != nil {
				_, err = c.Write([]byte(err.Error()))
				fmt.Println("Failed to prepare json response: marshal error")
				return
			}

			_, err = c.Write(respData)
		}
	}
}

func (s *UnixSocketServer) Connections() int {
	return s.connections
}

func AddTaskHandleFunc() {
	fmt.Println("Add Task")
}

func StartHandleFunc() {
	fmt.Println("Start")
}

func StopHandleFunc() {
	fmt.Println("Stop")
}

func PauseHandleFunc() {
	fmt.Println("Reset")
}

func ResetHandleFunc() {
	fmt.Println("Reset")
}

func StatusHandleFunc() {
	fmt.Println("Status")
}
