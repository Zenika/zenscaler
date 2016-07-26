package api

import (
	"net/http"
	"zscaler/core"

	"github.com/gin-gonic/gin"
)

// getScalers list all configured scalers
func getScalers(c *gin.Context) {
	var scalerNames = make([]string, 0)
	for k := range core.Config.Scalers {
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
		encoded, err := scaler.JSON()
		if err == nil {
			c.Data(http.StatusOK, "application/json", encoded)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func createScaler(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func patchScaler(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteScaler(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
