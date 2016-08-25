package api

import (
	"net/http"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/rule"

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
	var floatRuleBuilder FloatValueBuilder
	err := c.BindJSON(&floatRuleBuilder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON object cannot be parsed: " + err.Error(),
		})
		return
	}
	fvRule, err := floatRuleBuilder.Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON failed validation check: " + err.Error(),
		})
		return
	}
	// TODO check data race
	if _, exist := core.Config.Rules[fvRule.RuleName]; exist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "A rule named " + fvRule.RuleName + " already exist",
		})
		return
	}
	core.Config.Rules[fvRule.RuleName] = fvRule
	go rule.Watcher(core.Config.Errchan, fvRule)

	c.JSON(http.StatusCreated, gin.H{
		"rule": fvRule.RuleName,
	})
}

func patchRule(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Not implemented yet")
}

func deleteRule(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Not implemented yet")
}
