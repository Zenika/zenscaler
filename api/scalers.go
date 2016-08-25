package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zenika/zenscaler/core"
	"github.com/Zenika/zenscaler/core/scaler"
	"github.com/Zenika/zenscaler/core/types"

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

// ScalerBuilder contains all the information needed to build a scaler
type ScalerBuilder struct {
	Type string          `json:"type"`
	Name string          `json:"name"`
	Args json.RawMessage `json:"args"`
}

// Build validate inputed data and return a specific scaler wrapped in interface
func (sb *ScalerBuilder) Build() (types.Scaler, error) {
	if sb.Name == "" {
		return nil, fmt.Errorf("missing `name` field")
	}
	switch sb.Type {
	case "docker-compose-cmd":
		var dcs scaler.ComposeCmdScaler
		err := json.Unmarshal(sb.Args, &dcs)
		if err != nil {
			return nil, err
		}
		return &dcs, nil
	case "docker-service":
		var dss scaler.ServiceScaler
		err := json.Unmarshal(sb.Args, &dss)
		if err != nil {
			return nil, err
		}
		return &dss, nil
	default:
		return nil, fmt.Errorf("unknow scaler type")
	}
}

func createScaler(c *gin.Context) {
	var sb ScalerBuilder
	err := c.BindJSON(&sb)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON object cannot be parsed: " + err.Error(),
		})
		return
	}
	resultingScaler, err := sb.Build()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON failed validation check: " + err.Error(),
		})
		return
	}
	// TODO check data race
	if _, exist := core.Config.Scalers[sb.Name]; exist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "A scaler named" + sb.Name + "already exist",
		})
		return
	}
	core.Config.Scalers[sb.Name] = resultingScaler

	c.JSON(http.StatusCreated, gin.H{
		"scaler": sb.Name,
	})
}

func patchScaler(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Not implemented yet")
}

func deleteScaler(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Not implemented yet")
}
