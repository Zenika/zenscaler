package cmd

import (
	"github.com/Zenika/zenscaler/api"
	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/rule"
	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StartCmd parse configuration file and launch the scaler
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		setConfigPath()
		var err error
		core.Config, err = parseConfig()
		if err != nil {
			log.Fatalf("Fatal error reading config file: %s \n", err)
		}
		go api.Start()
		initialize()
	},
}

// Initialize core module
func initialize() {
	c := core.Config
	c.Errchan = make(chan error, 5)
	loop()
}

func loop() {
	c := core.Config
	log.Debug("Enter control loop...")
	// lanch a watcher on each rule
	for _, r := range c.Rules {
		go rule.Watcher(c.Errchan, r)
	}
	// watch for errors
	for {
		err := <-c.Errchan
		if err != nil {
			log.Errorf("%s", err)
		}
	}
}
