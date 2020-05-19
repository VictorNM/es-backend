package event

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"sync"
	"testing"
)

type UserRegistered struct {
	UserID int
}

func TestBus(t *testing.T) {
	t.Run("1 event", func(t *testing.T) {
		tests := []struct {
			nPublisher int
			nListener  int
		}{
			{1, 0},
			{1, 1},
			{1, 100},

			{0, 1},
			{1, 1},
			{100, 1},
		}

		for _, test := range tests {
			name := fmt.Sprintf("nPublisher = %d, nListener = %d", test.nPublisher, test.nListener)

			t.Run(name, func(t *testing.T) {
				b := NewBus()
				wg := &sync.WaitGroup{}
				wg.Add(test.nPublisher * test.nListener)

				ls := make([]*listener, test.nListener)
				for i := 0; i < test.nListener; i++ {
					ls[i] = newListener()
					ls[i].doSubscribe(b, UserRegistered{}, wg)
				}

				published := UserRegistered{UserID: 1}

				for i := 0; i < test.nPublisher; i++ {
					go b.Publish(published)
				}
				wg.Wait()

				for _, l := range ls {
					l.assertReceived(t, published, test.nPublisher)
				}
			})
		}
	})
}

type listener struct {
	c chan interface{}

	receivedEvents []interface{}
}

func newListener() *listener {
	return &listener{
		c: make(chan interface{}),
	}
}

func (l *listener) doSubscribe(b *Bus, e interface{}, wg *sync.WaitGroup) {
	b.Subscribe(e, l.c)
	go func() {
		for {
			select {
			case received := <-l.c:
				l.receivedEvents = append(l.receivedEvents, received)
				wg.Done()
			}
		}
	}()
}

func (l *listener) assertReceived(t *testing.T, e interface{}, times int) {
	assert.Equal(t, times, len(l.receivedEvents))
	for _, received := range l.receivedEvents {
		assert.Equal(t, e, received)
	}
}
