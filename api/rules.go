package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRule(c *gin.Context) {
	name := c.Param("name")
	// does this rule exist ?

	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func patchRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
