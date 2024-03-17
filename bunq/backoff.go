package bunq

import (
	"time"
)

const (
	Stop time.Duration = -1
)

// BackOff is a backoff policy for retrying an operation.
type backoff struct {
	// initialInterval is the initial interval for the first retry and upon the multiplication takes place.
	initialInterval time.Duration
	// maxInterval is the upper bound of backoff delay; the duration to sleep will not exceed this value.
	maxInterval time.Duration
	// maxTime is the maximum of elapsed time, before the backoff stops. If maxTime is zero, it means no maximum.
	maxTime time.Duration
	// multiplier is the factor to multiply the backoff interval.
	multiplier float64

	// State
	nextInterval time.Duration
	startTime    time.Time
	currentTry   int
}

func newBackoff(initialInterval, maxInterval, maxTime time.Duration, multiplier float64) backoff {
	if multiplier <= 1 {
		multiplier = 1.1
	}
	return backoff{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxTime:         maxTime,
		multiplier:      multiplier,
		nextInterval:    initialInterval,
	}
}

func NewDefaultBackoff() backoff {
	return newBackoff(1*time.Second, 4*time.Second, 12*time.Second, 1.25)
}

func (b *backoff) increaseInterval() {
	// Multiply the interval with the multiplier.
	interval := time.Duration(float64(b.nextInterval.Nanoseconds()) * b.multiplier)

	// Make sure the interval doesn't exceed the maxInterval.
	if interval > b.maxInterval {
		interval = b.maxInterval
	}

	b.nextInterval = interval
	b.currentTry++
}

// Try returns the current try, starting from 1 on the first try.
// Try is incremented after each call to NextLimit.
func (b *backoff) Try() int {
	return b.currentTry
}

// NextLimit returns the duration to wait before retrying the operation, or Stop to indicate that no more retries should be made.
// Stop will only be returned when the maxTime is set and the next interval would exceed the maxTime.
// Otherwise, the next interval is calculated and returned, unless it exceeds the maxInterval, in which case the maxInterval is returned.
func (b *backoff) NextLimit() (interval time.Duration) {
	// Set the initial time, so we can check if the maxTime is reached
	if b.maxTime > 0 && b.startTime.IsZero() {
		b.startTime = time.Now()
	}

	// Set the return value
	interval = b.nextInterval

	// Check whether if sleeping for the interval would exceed the maxTime.
	if b.maxTime > 0 && time.Since(b.startTime)+interval >= b.maxTime {
		interval = Stop
		//b.nextInterval = Stop
		return
	}

	// Increase the interval for the next run.
	b.increaseInterval()
	return
}

// Reset the interval back to the initial retry interval and restarts the timer.
// Reset can both be called after an operation, or right after one,
// as the time elapsed will be reset, but only started again when NextLimit is called.
func (b *backoff) Reset() {
	b.startTime = time.Time{}
	b.nextInterval = b.initialInterval
	b.currentTry = 0
}
