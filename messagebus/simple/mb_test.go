package simple_mb

import (
	"fmt"
	"sync"
	"testing"
)

func TestMb(t *testing.T)  {
	queueSize := 100
	bus := New(queueSize)


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


	bus.Publish("test", true, false)
	wg.Wait()
}
