package usecase

type customErr uint8

func (e customErr) Error() string {
	return errorsEnum[e]
}

const (
	// pomodoro errors
	MAX_COUNT customErr = iota
	UNABLE_TO_RESTART

	// tasks errors
	WRONG_ID customErr = iota
)

var errorsEnum = map[customErr]string{
	// pomodoro errors
	MAX_COUNT:         "maximum number of tomatoes",
	UNABLE_TO_RESTART: "cannot be restarted, task completed or already running",
	// tasks errors
	WRONG_ID: "invalid task id",
}
