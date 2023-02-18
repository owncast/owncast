package main

import (
	"fmt"

	ee "github.com/gravestench/eventemitter"
)

const (
	EventX = "x"
	EventY = "y"
	EventZ = "z"
)

func main() {
	bus := ee.New()

	NewService(bus, []string{"A", "B"}, "AB")
	NewService(bus, []string{"B", "C"}, "BC")
	NewService(bus, []string{"C", "D"}, "CD")

	fmt.Printf("sending [A, foo] ...\n")
	bus.Emit("A", "foo")
	fmt.Printf("\n")

	fmt.Printf("sending [B, bar] ...\n")
	bus.Emit("B", "bar")
	fmt.Printf("\n")

	fmt.Printf("sending [C, baz] ...\n")
	bus.Emit("C", "baz")
	fmt.Printf("\n")

	fmt.Printf("sending [d, qux] ...\n")
	bus.Emit("D", "qux")
}

type serviceCommon struct {
	Name      string
	bus       *ee.EventEmitter
	listenFor []string
}

func (s *serviceCommon) Attach(bus *ee.EventEmitter) {
	s.bus = bus
}

func NewService(bus *ee.EventEmitter, listenFor []string, name string) *Service {
	s := &Service{}

	s.Name = name
	s.bus = bus
	s.listenFor = listenFor

	s.init()

	return s
}

type Service struct {
	serviceCommon
}

func (s *Service) init() {
	for _, msg := range s.listenFor {
		event := msg
		s.bus.On(event, func(args ...any) {
			if len(args) > 0 {
				fmt.Printf("Service %s: [msg,data] => [%v, %v]\n", s.Name, event, args[0])
				return
			}

			fmt.Printf("Service %s: msg => %v\n", s.Name, event)
		})
	}
}
