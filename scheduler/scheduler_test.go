package scheduler

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestScheduler(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Start", func(t *testing.T) {
		var events []bool
		collector := make(chan bool)

		scheduler := New(time.Microsecond, func(innerStopped bool) {
			t.Logf("cb %v", innerStopped)
			collector <- innerStopped
			if innerStopped {
				close(collector)
			}
		})

		t.Log("Starting")
		scheduler.Start()

		// Wait for at least 2 invocations
		events = append(events, <-collector)
		events = append(events, <-collector)

		t.Log("Stopping")
		scheduler.Stop()

		t.Log("Waiting for finish ...")
		for ev := range collector {
			events = append(events, ev)
		}

		Expect(events[:len(events)-2]).To(SatisfyAll(BeFalse()))
		Expect(events[len(events)-1]).To(BeTrue())

		t.Log("Done")
	})
}
