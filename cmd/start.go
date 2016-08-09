package cmd

import (
	"github.com/Zenika/zscaler/api"
	"github.com/Zenika/zscaler/core"
	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// StartCmd parse configuration file and launch the scaler
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start autoscaler",
	Run: func(cmd *cobra.Command, args []string) {
		if Debug {
			log.SetLevel(log.DebugLevel)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		var err error
		core.Config, err = parseConfig()
		if err != nil {
			log.Fatalf("Fatal error reading config file: %s \n", err)
		}
		go api.Start()
		core.Config.Initialize()
	},
}
