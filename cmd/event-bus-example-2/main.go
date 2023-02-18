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

	NewServiceX(bus) // listen for EventX, emit EventY
	NewServiceY(bus) // listen for EventY, emit EventZ
	NewServiceZ(bus) // listen for EventZ

	bus.Emit(EventX, "foo")
}

type Service struct {
	serviceCommon
}

type serviceCommon struct {
	Name      string
	bus       *ee.EventEmitter
	listenFor []string
}

func NewServiceX(bus *ee.EventEmitter) *Service {
	s := &Service{}

	s.Name = "Service X"
	s.bus = bus

	s.bus.On(EventX, func(args ...any) {
		fmt.Printf("%s: message recieved\n", s.Name)
		s.bus.Emit(EventY, "bar")
	})

	return s
}

func NewServiceY(bus *ee.EventEmitter) *Service {
	s := &Service{}

	s.Name = "Service Y"
	s.bus = bus

	s.bus.On(EventY, func(args ...any) {
		fmt.Printf("%s: message recieved\n", s.Name)
		s.bus.Emit(EventZ, "baz")
	})

	return s
}

func NewServiceZ(bus *ee.EventEmitter) *Service {
	s := &Service{}

	s.Name = "Service Z"
	s.bus = bus

	s.bus.On(EventZ, func(args ...any) {
		fmt.Printf("%s: message recieved\n", s.Name)
	})

	return s
}
