package rule

import (
	"time"
	"zscaler/core/probe"
)

// A Rule must be able to perform a check
type Rule interface {
	Check() error // performe a check on the target and act if needed
	CheckInterval() time.Duration
}

// Watcher check periodically the rule and report back errors
// TODO channel back to kill if needed
func Watcher(c chan error, r Rule) {
	for {
		err := r.Check()
		if err != nil {
			c <- err
			return
		}
		time.Sleep(r.CheckInterval())
	}
}

// DefaultRule provide a basic implementation
type DefaultRule struct {
	Target Service
	Probe  probe.Probe
}

// Service describes the object to scale
type Service struct {
	Name  string
	Scale Scaler
}

// Check the mock probe
func (r DefaultRule) Check() error {
	if r.Probe.Value() > 0.75 {
		r.Target.Scale.Up()
	}
	if r.Probe.Value() < 0.25 {
		r.Target.Scale.Down()
	}
	return nil
}
