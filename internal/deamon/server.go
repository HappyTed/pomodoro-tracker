package deamon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"pomodoro.tracker/internal/entities"
)

// Unix Socker Server
type Server struct {
	lst            net.Listener
	MaxBufSize     int64 // kb, по дефолту можно сделать 128
	MaxConnections int
	ConnectionsNum int
	mx             sync.Mutex
}

func New(socketPath string, bufSize int64, maxConnections int) (*Server, error) {
	var s = &Server{MaxBufSize: bufSize, MaxConnections: maxConnections}

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
	fmt.Println("Try run unix socket domain server")
	fmt.Printf("Server run! Waiting for connection: %s\n", s.lst.Addr().String())

	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				fmt.Println("Current connections num:", s.ConnectionsNum, "/ MaxConn:", s.MaxConnections, ". Server listen:", s.lst.Addr().String())
				time.Sleep(1 * time.Second)
			}
		}

	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server stopped by context")
			s.lst.Close()
			return
		default:
			if s.ConnectionsNum < s.MaxConnections {
				nc, err := s.lst.Accept() // Это ожидающая операция
				// и к этому момоенту мог сработать контекст, получается гонка
				if err != nil {
					fmt.Println("Accept error:", err)
					continue
				}
				fmt.Println("New connection!")
				s.mx.Lock()
				s.ConnectionsNum++
				s.mx.Unlock()
				go s.handleCommand(ctx, nc)
			}

		}
	}

}

// работает с подключением до обрыва соединения
func (s *Server) handleCommand(ctx context.Context, c net.Conn) {
	defer func() {
		s.mx.Lock()
		s.ConnectionsNum--
		s.mx.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Server stopped by context")
			return
		default:
			buf := make([]byte, s.MaxBufSize)
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

			req := &entities.Request{}
			resp := &entities.Response{
				Status: entities.OK,
			}
			err = json.Unmarshal(data, req)
			if err != nil {
				fmt.Println("Marshal read data Error:", err)
				resp = &entities.Response{
					Status:  entities.ERROR,
					Message: err.Error(),
				}
			}

			// TODO: где-то здесь вызывать нужны обработчик
			if cmd, ok := entities.Commands[req.Cmd]; ok {
				switch cmd {
				case entities.ADD:
					AddTaskHandleFunc()
				case entities.START:
					StartHandleFunc()
				case entities.STOP:
					StopHandleFunc()
				case entities.PAUSE:
					PauseHandleFunc()
				case entities.RESET:
					ResetHandleFunc()
				case entities.STATUS:
					StatusHandleFunc()
				}
			} else {
				resp = &entities.Response{
					Status:  entities.ERROR,
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
