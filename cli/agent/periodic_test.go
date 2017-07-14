package agent

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"
	"speedland.net/service/go/src/lib/wcg"
)

// SimplePeriodicJob is a periodic job implementation that count up the # of task executions per 100 microseconds.
// it will stop after the count reached to the threashold (`until` field) or `stopped` is flagged.
type SimplePeriodicJob struct {
	c          int
	until      int
	panicAfter int
}

func (sr *SimplePeriodicJob) RunOnce(ctx context.Context) error {
	sr.c++
	if sr.panicAfter > 0 && sr.c > sr.panicAfter {
		panic(fmt.Errorf("panic after %d", sr.panicAfter))
	}
	return nil
}

func (sr *SimplePeriodicJob) ShouldRun(ctx context.Context) bool {
	if sr.until != 0 {
		return sr.c < sr.until
	}
	return true
}

func Test_Periodic_StopBySignal(t *testing.T) {
	assert := wcg.NewAssert(t)
	ctx := context.Background()
	job := &SimplePeriodicJob{}
	p := NewPeriodic(job, 100*time.Millisecond)
	p.Start(ctx)
	time.Sleep(1 * time.Second)
	p.Stop(ctx)

	assert.EqInt(10, job.c)
	assert.OK(!p.IsRunning())
}

func Test_Periodic_StopByTaskCondition(t *testing.T) {
	assert := wcg.NewAssert(t)
	ctx := context.Background()

	job := &SimplePeriodicJob{}
	job.until = 1
	p := NewPeriodic(job, 100*time.Millisecond)
	p.Start(ctx)
	time.Sleep(1 * time.Second)

	assert.EqInt(1, job.c)
	assert.OK(!p.IsRunning())

	// it's ok to call Stop() method only once (and usually should do this to release the resources)
	p.Stop(ctx)
}

func Test_Periodic_ContinueToRunAfterPanic(t *testing.T) {
	assert := wcg.NewAssert(t)
	ctx := context.Background()
	job := &SimplePeriodicJob{}
	job.panicAfter = 5
	p := NewPeriodic(job, 100*time.Millisecond)
	p.Start(ctx)
	time.Sleep(1 * time.Second)
	p.Stop(ctx)

	assert.EqInt(10, job.c)
	assert.OK(!p.IsRunning())
}
