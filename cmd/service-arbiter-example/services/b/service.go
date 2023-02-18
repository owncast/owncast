package b

import (
	"fmt"
)

type arbiter interface {
	CallServiceA()
	CallServiceC()
}

type Service struct {
	arbiter arbiter
}

func (s *Service) Call() {
	fmt.Println("Service B has been called")
	s.arbiter.CallServiceC()
}

func New(a arbiter) *Service {
	return &Service{
		arbiter: a,
	}
}
