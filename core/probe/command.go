package probe

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Command probe allow the execution of any file or script and parse the output
type Command struct {
	Cmd string
}

// Name of the probe
func (cp Command) Name() string {
	return "Command probe [" + cp.Cmd + "]"
}

// Value report the parsed output of the commmand or script as float64
func (cp Command) Value() float64 {
	output, err := cp.newCommand().Output()
	if err != nil {
		log.Errorf("Cannot probe [%s]: %s", cp.Cmd, err)
		return 0
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		log.Errorf("Cannot parse [%s] output to float: %s", cp.Cmd, err)
		return 0
	}
	return val
}

// NewCommand from string cut executable and args
func (cp Command) newCommand() *exec.Cmd {
	splitted := strings.Split(cp.Cmd, " ")
	return exec.Command(splitted[0], splitted[1:]...)
}
