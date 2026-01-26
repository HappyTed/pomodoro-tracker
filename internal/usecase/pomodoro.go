package usecase

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type State int8

const (
	STOPPED State = iota
	ACTIVE
	PAUSE
	DONE
)

const (
	WORK_TIME      = time.Minute * 25
	BREAK_TIME     = time.Minute * 5
	BIG_BREAK_TIME = time.Minute * 15
)

type Pomodoro struct {
	WaitGroup sync.WaitGroup
	mu        sync.Mutex
	id        uuid.UUID
	state     State
	counter   uint8
	timer     time.Duration
	controlCh <-chan State
}

func NewPomodoro(controlChannel <-chan State) (*Pomodoro, error) {
	return &Pomodoro{
		state:     STOPPED,
		counter:   0,
		timer:     WORK_TIME,
		id:        uuid.New(),
		controlCh: controlChannel,
	}, nil
}

func (p *Pomodoro) Run() error {
	p.state = ACTIVE
	return nil
}

func (p *Pomodoro) Pause() error {
	p.state = PAUSE
	return nil
}

func (p *Pomodoro) Stop() error {
	p.state = STOPPED
	return nil
}

func (p *Pomodoro) Next() error {
	if p.counter >= 254 {
		return MAX_COUNT
	}
	p.counter += 1
	if p.counter%4 == 0 {
		p.timer = BIG_BREAK_TIME
	} else if p.counter%2 == 0 {
		p.timer = BREAK_TIME
	} else {
		p.timer = WORK_TIME
	}

	p.Run()

	return nil
}

func (p *Pomodoro) ID() string {
	return p.id.String()
}

func switcher(p *Pomodoro) error {
	switch p.state {
	case ACTIVE, DONE:
		return UNABLE_TO_RESTART
	case STOPPED, PAUSE:
		err := p.Run()
		return err
	default:
		return nil
	}
}
