package backoff

import (
	"context"
	"fmt"
	"testing"
)

func TestBackoff(t *testing.T) {
	ctx := context.Background()
	bo := NewBackoffer(1000, ctx)
	var err error
	err = fmt.Errorf("test error")

	for i := 0; i < 3; i++ {
		if err != nil {
			if err := bo.Backoff(boTiKVRPC, err); err != nil {
				fmt.Println(err)
			}
		}

	}
}
