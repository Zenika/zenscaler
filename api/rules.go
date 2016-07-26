package api

import (
	"net/http"
	"zscaler/core"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func getRules(c *gin.Context) {
	var ruleNames = make([]string, 0)
	for k := range core.Config.Rules {
		ruleNames = append(ruleNames, k)
	}
	c.JSON(http.StatusOK, gin.H{
		"rules": ruleNames,
	})
}

func getRule(c *gin.Context) {
	name := c.Param("name")
	// does this rule exist ?
	if rule, ok := core.Config.Rules[name]; ok {
		encoded, err := rule.JSON()
		if err == nil {
			c.Data(http.StatusOK, "application/json", encoded)
			return
		}
		log.Errorf("Encode error %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": name + " does not encode to JSON",
		})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": name + " not found",
	})
}

func createRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func patchRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}

func deleteRule(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, "Not implemented yet")
}
