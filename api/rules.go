package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getRules(c *gin.Context) {
	name := c.Param("name")
	// does this rule exist ?

	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func patchRules(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteRules(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
