package scaler

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

const defaultTimeout = 3

// ComposeScaler hold configuration for a libcompose scaler
type ComposeScaler struct {
	ServiceName     string `json:"service"` // name of the service to scale
	ConfigFile      string `json:"config"`
	ProjectName     string `json:"project"` // usually the config file directory name
	UpperCountLimit uint64 `json:"upperCountLimit"`
	LowerCountLimit uint64 `json:"lowerCountLimit"`
	project         project.APIProject
}

// NewComposeScaler create a scaler managed by libcompose
//
// This assume that the project has been started beforehand
// with docker-compose python cli !
func NewComposeScaler(name, projectName, configFilePath string) (*ComposeScaler, error) {
	cs := &ComposeScaler{
		ServiceName:     name,
		ProjectName:     projectName,
		ConfigFile:      configFilePath,
		UpperCountLimit: 0, // default to unlimited
		LowerCountLimit: 1, // default to one, ensuring service avaibility
	}

	ctx := &docker.Context{
		Context: project.Context{
			ComposeFiles: []string{configFilePath},
			ProjectName:  projectName,
		}}

	var err error
	cs.project, err = docker.NewProject(ctx, &config.ParseOptions{})
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// Describe scaler
func (cs *ComposeScaler) Describe() string {
	return "Docker libcompose scaler"
}

// JSON encoding
func (cs *ComposeScaler) JSON() ([]byte, error) {
	encoded, err := json.Marshal(*cs)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// Up using libcompose API
func (cs *ComposeScaler) Up() error {
	target := cs.countContainers() + 1
	if target < cs.UpperCountLimit && cs.UpperCountLimit != 0 {
		return cs.scale(target)
	}
	cs.getLogger().Debug("cannot scale up: maximum count achieved")
	return nil
}

// Down using libcompose API
func (cs *ComposeScaler) Down() error {
	target := cs.countContainers() - 1
	if target > cs.LowerCountLimit {
		return cs.scale(target)
	}
	cs.getLogger().Debug("cannot scale down: munimum count achieved")
	return nil
}

// scaler with libcompose to the targeted value
func (cs *ComposeScaler) scale(target uint64) error {
	servicesScale := map[string]int{cs.ServiceName: int(target)}
	err := cs.project.Up(context.Background(), options.Up{})
	if err != nil {
		return err
	}
	err = cs.project.Scale(context.Background(), defaultTimeout, servicesScale)
	return err
}

// countContainers running the specified service in this project
func (cs *ComposeScaler) countContainers() uint64 {
	inf, err := cs.project.Ps(context.Background(), false)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	fmt.Printf("%v", inf)
	return 0
}

func (cs *ComposeScaler) getLogger() *log.Entry {
	return log.WithFields(log.Fields{
		"service": cs.ServiceName,
		"project": cs.ProjectName,
	})
}
