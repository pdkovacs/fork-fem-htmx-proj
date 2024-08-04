package incrementor

import (
	"testing"
	"time"
)

func TestStartIncrementing(t *testing.T) {
	stopResetIncrementing()

	incrementBy := 100
	intervalMs := 100
	incrementCount := 7
	StartIncrementing(incrementBy, intervalMs)
	time.Sleep(time.Duration(incrementCount*intervalMs*1e6 + intervalMs*1e6/2))
	want := incrementCount * incrementBy
	if consumed != want {
		t.Fatalf("Expected %d, got %d", want, consumed)
	}
}

func TestStartWhileAlreadyIncrementing(t *testing.T) {
	stopResetIncrementing()

	incrementBy := 100
	intervalMs := 100
	incrementCount := 7
	startTime := time.Now()
	StartIncrementing(incrementBy, intervalMs)

	incrementByBusily := incrementBy * 1000
	time.Sleep(time.Duration(incrementCount*intervalMs*1e6/2 + intervalMs*1e6/2))
	busyStart := time.Now()
	StartIncrementing(incrementByBusily, 1)
	time.Sleep(time.Duration(incrementCount * intervalMs * 1e6))

	sinceBusyStartMs := int(time.Since(busyStart) / 1e6)
	if sinceBusyStartMs < intervalMs*3 {
		t.Fatalf("elapsed time since the second start of incrementing (%d ms) is less then the tripple interval (%d ms): test will not be reliable", sinceBusyStartMs, intervalMs*3)
	}

	elapsedTimeInMs := int(time.Since(startTime)) / 1e6
	actualIntervalCount := elapsedTimeInMs / intervalMs
	wantLessThan := (actualIntervalCount + 1) * incrementBy
	if consumed >= wantLessThan {
		t.Fatalf("Expected less than %d, got %d after %d ms \n", wantLessThan, consumed, elapsedTimeInMs)
	}
}

func TestStartSuspendIncrementing(t *testing.T) {
	stopResetIncrementing()

	incrementBy := 100
	intervalMs := 100
	incrementCount := 7
	StartIncrementing(incrementBy, intervalMs)

	letItRunForMs := incrementCount * intervalMs / 2
	time.Sleep(time.Duration(letItRunForMs * 1e6))
	SuspendIncrementing()
	time.Sleep(time.Duration(incrementCount * intervalMs * 1e6 / 2))
	SuspendIncrementing()

	actualIntervalCount := letItRunForMs / intervalMs
	want := actualIntervalCount * incrementBy
	if consumed != want {
		t.Fatalf("Expected %d, got %d", want, consumed)
	}
}

func TestStartSuspendResumeIncrementing(t *testing.T) {
	stopResetIncrementing()

	incrementBy := 100
	intervalMs := 100
	incrementCount := 7
	letItRunForMs := incrementCount * intervalMs
	suspendTimeMs := int(float64(incrementCount)*1.5) * intervalMs

	startTime := time.Now()
	StartIncrementing(incrementBy, intervalMs)
	time.Sleep(time.Duration(letItRunForMs * 1e6 / 2))
	SuspendIncrementing()
	time.Sleep(time.Duration(suspendTimeMs * 1e6))
	SuspendIncrementing()
	StartIncrementing(incrementBy, intervalMs) // typically waits an extra interval
	time.Sleep(time.Duration(letItRunForMs * 1e6 / 2))
	SuspendIncrementing()
	endTime := time.Now()

	activeRunMs := int(endTime.Sub(startTime).Milliseconds()) - suspendTimeMs
	actualIntervalCount := activeRunMs / intervalMs
	accountForTheExtraIntervalOnResume := -1
	want := (actualIntervalCount + accountForTheExtraIntervalOnResume) * incrementBy
	if consumed != want {
		t.Fatalf("Expected %d, got %d", want, consumed)
	}
}
