package entities

type Command uint8

const (
	START Command = iota
	STOP
	PAUSE
	RESET
)

type Request struct {
	Cmd  Command  `json:"cmd"`
	Args []string `json:"args,omitempty"`
}

type Status uint

const (
	OK Status = iota
	ERROR
)

type Response struct {
	Status  Status   `json:"status"`
	Message *string  `json:"message,omitempty"`
	Data    []string `json:"data,omitempty"`
}
