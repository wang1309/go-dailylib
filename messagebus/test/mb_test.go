package test

import (
	"fmt"
	messagebus "github.com/vardius/message-bus"
	"sync"
	"testing"
)

func TestMb(t *testing.T)  {
	queueSize := 100
	bus := messagebus.New(queueSize)


	var wg sync.WaitGroup
	wg.Add(2)

	_ = bus.Subscribe("test", func(v bool) {
		defer wg.Done()
		fmt.Println(v)
	})

	_ = bus.Subscribe("test", func(v bool) {
		defer wg.Done()
		fmt.Println(v)
	})


	bus.Publish("test", true)
	wg.Wait()
}
