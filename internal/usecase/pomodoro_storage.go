package usecase

type PomodoroManager struct {
	tomatos        map[string]*Pomodoro
	controlChannel chan State
}

func NewPomodoroManager() (*PomodoroManager, error) {
	return &PomodoroManager{
		tomatos: make(map[string]*Pomodoro),
	}, nil
}

func (tm *PomodoroManager) List() ([]*Pomodoro, error) {
	var result []*Pomodoro
	for _, p := range tm.tomatos {
		result = append(result, p)
	}
	return result, nil
}

func (tm *PomodoroManager) Add() (string, error) {
	p, err := NewPomodoro(tm.controlChannel)
	if err != nil {
		return "", err
	}
	uid := p.ID()
	tm.tomatos[uid] = p

	return uid, err
}
