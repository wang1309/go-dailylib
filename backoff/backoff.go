package backoff

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	// NoJitter makes the backoff sequence strict exponential.
	NoJitter = 1 + iota
	// FullJitter applies random factors to strict exponential.
	FullJitter
	// EqualJitter is also randomized, but prevents very short sleeps.
	EqualJitter
	// DecorrJitter increases the maximum jitter based on the last random value.
	DecorrJitter
)

// NewBackoffFn creates a backoff func which implements exponential backoff with
// optional jitters.
// See http://www.awsarchitectureblog.com/2015/03/backoff.html
func NewBackoffFn(base, cap, jitter int) func(ctx context.Context) int {
	attempts := 0
	lastSleep := base
	return func(ctx context.Context) int {
		var sleep int
		switch jitter {
		case NoJitter:
			sleep = expo(base, cap, attempts)
		case FullJitter:
			v := expo(base, cap, attempts)
			sleep = rand.Intn(v)
		case EqualJitter:
			v := expo(base, cap, attempts)
			sleep = v/2 + rand.Intn(v/2)
		case DecorrJitter:
			sleep = int(math.Min(float64(cap), float64(base+rand.Intn(lastSleep*3-base))))
		}

		select {
		case <-time.After(time.Duration(sleep) * time.Millisecond):
			fmt.Println("ticker: ", sleep)
		case <-ctx.Done():
		}

		attempts++
		lastSleep = sleep
		return lastSleep
	}
}

func expo(base, cap, n int) int {
	return int(math.Min(float64(cap), float64(base)*math.Pow(2.0, float64(n))))
}

type backoffType int

// Back off types.
const (
	boTiKVRPC backoffType = iota
	BoTxnLock
	boTxnLockFast
	boPDRPC
	BoRegionMiss
	boServerBusy
)

func (t backoffType) createFn() func(context.Context) int {
	switch t {
	case boTiKVRPC:
		return NewBackoffFn(100, 2000, EqualJitter)
	case BoTxnLock:
		return NewBackoffFn(200, 3000, EqualJitter)
	case boTxnLockFast:
		return NewBackoffFn(100, 3000, EqualJitter)
	case boPDRPC:
		return NewBackoffFn(500, 3000, EqualJitter)
	case BoRegionMiss:
		return NewBackoffFn(100, 500, NoJitter)
	case boServerBusy:
		return NewBackoffFn(2000, 10000, EqualJitter)
	}
	return nil
}

func (t backoffType) String() string {
	switch t {
	case boTiKVRPC:
		return "tikvRPC"
	case BoTxnLock:
		return "txnLock"
	case boTxnLockFast:
		return "txnLockFast"
	case boPDRPC:
		return "pdRPC"
	case BoRegionMiss:
		return "regionMiss"
	case boServerBusy:
		return "serverBusy"
	}
	return ""
}

// Maximum total sleep time(in ms) for kv/cop commands.
const (
	copBuildTaskMaxBackoff  = 5000
	tsoMaxBackoff           = 5000
	scannerNextMaxBackoff   = 20000
	batchGetMaxBackoff      = 20000
	copNextMaxBackoff       = 20000
	getMaxBackoff           = 20000
	prewriteMaxBackoff      = 20000
	cleanupMaxBackoff       = 20000
	GcOneRegionMaxBackoff   = 20000
	GcResolveLockMaxBackoff = 100000
	GcDeleteRangeMaxBackoff = 100000
	rawkvMaxBackoff         = 20000
	splitRegionBackoff      = 20000
)

var commitMaxBackoff = 20000

// Backoffer is a utility for retrying queries.
type Backoffer struct {
	context.Context

	fn         map[backoffType]func(context.Context) int
	maxSleep   int
	totalSleep int
	errors     []error
	types      []backoffType
}

// NewBackoffer creates a Backoffer with maximum sleep time(in ms).
func NewBackoffer(maxSleep int, ctx context.Context) *Backoffer {
	return &Backoffer{
		Context:  ctx,
		maxSleep: maxSleep,
	}
}

// Backoff sleeps a while base on the backoffType and records the error message.
// It returns a retryable error if total sleep time exceeds maxSleep.
func (b *Backoffer) Backoff(typ backoffType, err error) error {
	select {
	case <-b.Context.Done():
		return err
	default:
	}

	// Lazy initialize.
	if b.fn == nil {
		b.fn = make(map[backoffType]func(context.Context) int)
	}

	f, ok := b.fn[typ]
	if !ok {
		f = typ.createFn()
		b.fn[typ] = f
	}

	b.totalSleep += f(b)
	b.types = append(b.types, typ)

	b.errors = append(b.errors, err)
	if b.maxSleep > 0 && b.totalSleep >= b.maxSleep {
		errMsg := fmt.Sprintf("backoffer.maxSleep %dms is exceeded, errors:", b.maxSleep)
		for i, err := range b.errors {
			// Print only last 3 errors for non-DEBUG log levels.
			if i >= len(b.errors)-3 {
				errMsg += "\n" + err.Error()
			}
		}
		log.Println(errMsg)
		// Use the last backoff type to generate a MySQL error.
		return typ.TError()
	}
	return nil
}

func (t backoffType) TError() error {
	/*	switch t {
		case boTiKVRPC:
			return ErrTiKVServerTimeout
		case BoTxnLock, boTxnLockFast:
			return ErrResolveLockTimeout
		case boPDRPC:
			return ErrPDServerTimeout.GenByArgs(txnRetryableMark)
		case BoRegionMiss:
			return ErrRegionUnavailable
		case boServerBusy:
			return ErrTiKVServerBusy
		}*/
	return nil
}

// Clone creates a new Backoffer which keeps current Backoffer's sleep time and errors, and shares
// current Backoffer's context.
func (b *Backoffer) Clone() *Backoffer {
	return &Backoffer{
		Context:    b.Context,
		maxSleep:   b.maxSleep,
		totalSleep: b.totalSleep,
		errors:     b.errors,
	}
}

// Fork creates a new Backoffer which keeps current Backoffer's sleep time and errors, and holds
// a child context of current Backoffer's context.
func (b *Backoffer) Fork() (*Backoffer, context.CancelFunc) {
	ctx, cancel := context.WithCancel(b.Context)
	return &Backoffer{
		Context:    ctx,
		maxSleep:   b.maxSleep,
		totalSleep: b.totalSleep,
		errors:     b.errors,
	}, cancel
}