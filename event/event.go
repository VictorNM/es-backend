package event

import (
	"reflect"
)

type Bus struct {
	eventMap map[string][]chan<- interface{}
}

func (b *Bus) Publish(e interface{}) {
	name := reflect.TypeOf(e).Name()

	for _, c := range b.eventMap[name] {
		c <- e
	}
}

func (b *Bus) Subscribe(e interface{}, c chan<- interface{}) {
	name := reflect.TypeOf(e).Name()

	b.eventMap[name] = append(b.eventMap[name], c)
}

func NewBus() *Bus {
	return &Bus{
		eventMap: make(map[string][]chan<- interface{}),
	}
}