// Package agent provides the utilities to implement long running command line application
package agent

import (
	"golang.org/x/net/context"
)

// Agent is an interface to implement an agent feature
type Agent interface {
	// Start starts the agent
	Start(ctx context.Context) error
	// Stop stops the agent
	Stop(ctx context.Context) error
	// IsRunning returns the agent is running or not
	IsRunning() bool
}
