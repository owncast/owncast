package c

import (
	"fmt"
)

type arbiter interface {
	CallServiceA()
	CallServiceB()
}

type Service struct {
	arbiter arbiter
}

func (s *Service) Call() {
	fmt.Println("Service C has been called")
}

func New(a arbiter) *Service {
	return &Service{
		arbiter: a,
	}
}
