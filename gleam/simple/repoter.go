package simple

import (
	"context"
	"sync"
)

func (runner *gleamRunner) report(ctx context.Context, f func() error) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	finishedChan := make(chan bool, 1)

	var heartbeatWg sync.WaitGroup
	heartbeatWg.Add(1)
	defer heartbeatWg.Wait()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(finishedChan)
	}()

	select {
	case <-finishedChan:
		runner.reportStatus()
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
