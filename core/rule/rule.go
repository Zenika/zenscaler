package rule

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
	"zscaler/core/probe"
	"zscaler/core/scaler"

	log "github.com/Sirupsen/logrus"
)

// A Rule must be able to perform a check
type Rule interface {
	Check() error                 // performe a check on the target and act if needed
	CheckInterval() time.Duration // time to wait between each check
	JSON() ([]byte, error)        // return a JSON output
}

// Watcher check periodically the rule and report back errors
// TODO channel back to kill if needed
func Watcher(c chan error, r Rule) {
	for {
		err := r.Check()
		if err != nil {
			c <- fmt.Errorf("Error checking probe: %s", err)
			return
		}
		time.Sleep(r.CheckInterval())
	}
}

// FloatValue handler
type FloatValue struct {
	RuleName       string               `json:"rule"`
	ServiceName    string               `json:"service"`
	Scale          scaler.Scaler        `json:"-"`
	ScalerID       string               `json:"scaler"`
	Probe          probe.Probe          `json:"-"`
	ProbeID        string               `json:"probe"`
	RefreshRate    time.Duration        `json:"resfreshRate"`
	UpDefinition   string               `json:"up"`
	Up             func(v float64) bool `json:"-"`
	DownDefinition string               `json:"down"`
	Down           func(v float64) bool `json:"-"`
}

// Check the probe, UP and DOWN
func (r *FloatValue) Check() error {
	probe, err := r.Probe.Value()
	if err != nil {
		return err
	}
	log.Debugf("["+r.ServiceName+"] "+r.Probe.Name()+" at %.2f ", probe)
	if r.Up(probe) && r.Down(probe) {
		log.Warning("[" + r.ServiceName + "] try to scale up and down at the same time! (nothing done)")
		return nil
	}
	if r.Up(probe) {
		err := r.Scale.Up()
		if err != nil {
			log.Errorf("Error when scaling up: %s", err)
			return nil
		}
	}
	if r.Down(probe) {
		err := r.Scale.Down()
		if err != nil {
			log.Errorf("Error when scaling down: %s", err)
			return nil
		}
	}
	return nil
}

// CheckInterval return the time to wait between each check
func (r *FloatValue) CheckInterval() time.Duration {
	return r.RefreshRate
}

// JSON encode
func (r *FloatValue) JSON() ([]byte, error) {
	encoded, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Parse UpDefinition and DownDefinition directives to create matching functions
func (r *FloatValue) Parse() (err error) {
	r.Up, err = Decode(r.UpDefinition)
	if err != nil {
		return errors.New(r.ServiceName + fmt.Sprintf(": %v up", err))
	}
	r.Down, err = Decode(r.DownDefinition)
	if err != nil {
		return errors.New(r.ServiceName + fmt.Sprintf(": %v down", err))
	}
	return
}

// Decode a logical rule (ex. ">0.75")
func Decode(order string) (func(float64) bool, error) {
	// check syntax
	regex := regexp.MustCompile(`^[[:space:]]*([>|<|==|!=])[[:space:]]*([[:digit:]]*(?:\.[[:digit:]]*)?)$`)
	matches := regex.FindStringSubmatch(order)
	if len(matches) == 3 {
		value, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return nil, nil
		}
		switch matches[1] {
		case ">":
			return func(p float64) bool { return p > value }, nil
		case "<":
			return func(p float64) bool { return p < value }, nil
		case "==":
			return func(p float64) bool { return p == value }, nil
		case "!=":
			return func(p float64) bool { return p != value }, nil
		}
	}
	return func(p float64) bool { return false }, errors.New("Error decoding rule [" + order + "]")
}
