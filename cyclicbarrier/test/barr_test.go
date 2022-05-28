package test

import (
	"context"
	"fmt"
	"github.com/marusama/cyclicbarrier"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)


func checkBarrier(t *testing.T, b cyclicbarrier.CyclicBarrier, expectedParties, expectedNumberWaiting int, expectedIsBroken bool) {

	parties, numberWaiting := b.GetParties(), b.GetNumberWaiting()
	isBroken := b.IsBroken()

	if expectedParties >= 0 && parties != expectedParties {
		t.Error("barrier must have parties = ", expectedParties, ", but has ", parties)
	}
	if expectedNumberWaiting >= 0 && numberWaiting != expectedNumberWaiting {
		t.Error("barrier must have numberWaiting = ", expectedNumberWaiting, ", but has ", numberWaiting)
	}
	if isBroken != expectedIsBroken {
		t.Error("barrier must have isBroken = ", expectedIsBroken, ", but has ", isBroken)
	}
}

func TestAwaitOnce(t *testing.T) {
	n := 2 // goroutines count
	b := cyclicbarrier.New(n)
	ctx := context.Background()

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			fmt.Println("before")
			err := b.Await(ctx)
			fmt.Println("after")
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	checkBarrier(t, b, n, 0, false)
}


func TestAwaitOnceCtxDone(t *testing.T) {
	n := 100        // goroutines count
	b := cyclicbarrier.New(n + 1) // parties are more than goroutines count so all goroutines will wait
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	var deadlineCount, brokenBarrierCount int32

	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(num int) {
			time.Sleep(100*time.Millisecond)
			fmt.Println("before")
			err := b.Await(ctx)
			fmt.Println("after")
			if err == context.DeadlineExceeded {
				//fmt.Println("deadline exceed")
				atomic.AddInt32(&deadlineCount, 1)
			} else if err == cyclicbarrier.ErrBrokenBarrier {
				fmt.Println("broken barrier")
				atomic.AddInt32(&brokenBarrierCount, 1)
			} else {
				panic("must be context.DeadlineExceeded or ErrBrokenBarrier error")
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	checkBarrier(t, b, n+1, -1, true)
	if deadlineCount == 0 {
		t.Error("must be more than 0 context.DeadlineExceeded errors, but found", deadlineCount)
	}
	if deadlineCount+brokenBarrierCount != int32(n) {
		t.Error("must be exactly", n, "context.DeadlineExceeded and ErrBrokenBarrier errors, but found", deadlineCount+brokenBarrierCount)
	}
}
