package services

import (
	"errors"
	"os"
	"time"
)

type InputService struct{}

func NewInputService() *InputService {
	return &InputService{}
}

func (input *InputService) GetArgs() (map[int]time.Time, error) {
	today := time.Now()

	var arguments = map[int]time.Time{
		0: time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location()),
	}

	for i, arg := range os.Args[1:] {
		argDate, err := time.Parse("2006-01-02", arg)

		if err != nil {
			return nil, errors.New("date should be formatted as YYYY-MM-DD")
		}

		arguments[i] = time.Date(argDate.Year(), argDate.Month(), argDate.Day(), 0, 0, 0, 0, today.Location())

		if i == 1 && arguments[i].Before(arguments[0]) {
			return nil, errors.New("dates must be in ascending order")
		}
	}
	return arguments, nil
}
