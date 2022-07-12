package test

import (
	"fmt"
	"testing"
	"time"
)

func TestFmt(t *testing.T)  {
	// Make a big map.
	m := map[int]int{}
	for i := 0; i < 100000; i++ {
		m[i] = i
	}
	c := make(chan string)
	go func() {
		// Print the map.
		s := fmt.Sprintln(m)
		c <- s
	}()
	go func() {
		time.Sleep(1 * time.Millisecond)
		// Add an extra item to the map while iterating.
		m[-1] = -1
		c <- ""
	}()
	<-c
	<-c
}