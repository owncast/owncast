package a

import (
	"fmt"
)

type arbiter interface {
	CallServiceB()
	CallServiceC()
}

type Service struct {
	arbiter arbiter
}

func (s *Service) Call() {
	fmt.Println("Service A has been called")
	s.arbiter.CallServiceC()
}

func New(a arbiter) *Service {
	return &Service{
		arbiter: a,
	}
}
