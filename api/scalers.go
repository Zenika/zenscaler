package api

import (
	"net/http"
	"zscaler/core"

	"github.com/gin-gonic/gin"
)

// getScalers list all configured scalers
func getScalers(c *gin.Context) {
	var scalerNames = make([]string, 0)
	for k := range core.Config.Rules {
		scalerNames = append(scalerNames, k)
	}
	c.JSON(http.StatusOK, gin.H{
		"scalers": scalerNames,
	})
}

// getScaler give the configuration parameters of a scaler
func getScaler(c *gin.Context) {
	name := c.Param("name")
	// does this scaler exist ?
	if scaler, ok := core.Config.Scalers[name]; ok {
		c.JSON(http.StatusOK, scaler)
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func patchScaler(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteScaler(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
