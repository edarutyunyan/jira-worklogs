package services

type Services struct {
	Output OutputService
}

func NewServices(output OutputService) *Services {
	return &Services{
		Output: output,
	}
}
