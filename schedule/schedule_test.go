package schedule

import (
	"context"
	"fmt"
	"testing"
)

func TestSchedule(t *testing.T)  {
	schedule := NewFIFOScheduler()
	defer func() {
		schedule.Stop()
	}()

	schedule.Schedule(func(ctx context.Context) {
		fmt.Println("test schedule")
	})
}
