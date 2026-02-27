package entities

type Command uint8

const (
	ADD Command = iota
	START
	STOP
	PAUSE
	RESET
	STATUS
)

var Commands map[string]Command = map[string]Command{
	"ADD":    ADD,
	"START":  START,
	"STOP":   STOP,
	"PAUSE":  PAUSE,
	"RESET":  RESET,
	"STATUS": STATUS,
}

type Request struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args,omitempty"`
}

type Status uint

const (
	OK Status = iota
	ERROR
)

type Response struct {
	Status  Status   `json:"status"`
	Message string   `json:"message,omitempty"`
	Data    []string `json:"data,omitempty"`
}
