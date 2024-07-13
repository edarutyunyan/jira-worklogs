package services

type Services struct {
	Input  InputService
	Output OutputService
}

func NewServices(input InputService, output OutputService) *Services {
	return &Services{
		Input:  input,
		Output: output,
	}
}
