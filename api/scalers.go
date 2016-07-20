package api

import (
	"net/http"
	"zscaler/core"

	"github.com/gin-gonic/gin"
)

func getScalers(c *gin.Context) {
	name := c.Param("name")
	// does this scaler exist ?
	if scaler, ok := core.Config.Scalers[name]; ok {
		c.JSON(http.StatusOK, scaler)
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func patchScalers(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteScalers(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
