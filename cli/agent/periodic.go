package agent

import "time"
import "golang.org/x/net/context"

import "github.com/speedland/go/x/xlog"

// Periodic is an Agent implementation to run a job by periodic intervals.
type Periodic struct {
	Interval   time.Duration
	job        PeriodicJob
	jobContext context.Context
	cancelFunc context.CancelFunc
}

// PeriodicJob is a job interface for a periodic agent to run a job
type PeriodicJob interface {
	RunOnce(ctx context.Context) error
	ShouldRun(ctx context.Context) bool
}

func NewPeriodic(job PeriodicJob, interval time.Duration) *Periodic {
	return &Periodic{
		Interval: interval,
		job:      job,
	}
}

// Start implements Agent#Start()
func (p *Periodic) Start(ctx context.Context) error {
	p.jobContext, p.cancelFunc = context.WithCancel(ctx)
	go p.main(p.jobContext)
	return nil
}

// IsRunning implements Agent#IsRunning()
func (p *Periodic) IsRunning() bool {
	return p.jobContext != nil
}

// Stop implements Agent#Stop()
func (p *Periodic) Stop(ctx context.Context) error {
	if p.jobContext != nil {
		err := p.jobContext.Err()
		p.cancelFunc()
		p.cancelFunc = nil
		p.jobContext = nil
		return err
	}
	return nil
}

func (p *Periodic) main(ctx context.Context) {
	var logger = xlog.WithContext(ctx)
	defer func() {
		p.cancelFunc = nil
		p.jobContext = nil
	}()
	for {
		if p.job.ShouldRun(ctx) {
			go func() {
				defer func() {
					if x := recover(); x != nil {
						logger.Errorf("a periodic job fails with a panic: %v", x)
					}
				}()
				p.job.RunOnce(ctx)
			}()
		} else {
			return
		}
		select {
		case <-time.Tick(p.Interval):
			break
		case <-ctx.Done():
			// canceled
			return
		}
	}
}
