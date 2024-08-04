package incrementor

import (
	"fmt"
	"sync"
	"time"
)

var started bool = false
var stopRequested bool = false
var consumed int = 0
var timer *time.Timer
var timerMux *sync.Mutex = &sync.Mutex{}

func StartIncrementing(incrementBy int, intervalMs int) {
	timerMux.Lock()
	defer timerMux.Unlock()

	fmt.Printf("startIncrementing: started=%#v\n", started)
	if started {
		fmt.Printf("startIncrementing: already started; returning...\n")
		return
	}
	started = true
	stopRequested = false

	incrementAtIntevals(incrementBy, intervalMs)
}

func GetConsumed() int {
	return consumed
}

func SuspendIncrementing() {
	timerMux.Lock()
	defer timerMux.Unlock()

	if !started {
		return
	}
	started = false

	stopRequested = true

	if !timer.Stop() {
		<-timer.C
	}
}

func incrementAtIntevals(incrementBy int, intervalMs int) {

	timer = time.NewTimer(time.Duration(intervalMs * 1e6))

	go func() {
		<-timer.C

		timerMux.Lock()
		defer timerMux.Unlock()

		if !stopRequested {
			fmt.Printf("Just about to increment by %d...\n", incrementBy)
			consumed += incrementBy
			go func() {
				incrementAtIntevals(incrementBy, intervalMs)
			}()
		}
	}()
}

func stopResetIncrementing() {
	SuspendIncrementing()

	timerMux.Lock()
	defer timerMux.Unlock()
	consumed = 0
}
