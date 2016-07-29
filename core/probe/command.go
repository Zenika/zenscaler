package probe

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Command probe allow the execution of any file or script and parse the output
type Command struct {
	Cmd string `json:"cmd"`
}

// Name of the probe
func (cp Command) Name() string {
	return "Command probe [" + cp.Cmd + "]"
}

// Value report the parsed output of the commmand or script as float64
func (cp Command) Value() (float64, error) {
	output, err := cp.newCommand().Output()
	if err != nil {
		return 0, fmt.Errorf("Cannot probe [%s]: %s", cp.Cmd, err)
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse [%s] output to float: %s", cp.Cmd, err)
	}
	return val, nil
}

// NewCommand from string cut executable and args
func (cp Command) newCommand() *exec.Cmd {
	splitted := strings.Split(cp.Cmd, " ")
	return exec.Command(splitted[0], splitted[1:]...)
}
