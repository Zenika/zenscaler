package api

import "github.com/gin-gonic/gin"

// Start the API listener
func Start() error {
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.GET("/scalers", getScalers)
	v1.PATCH("/scalers/:name", patchScalers)
	v1.DELETE("/scalers/:name", deleteScalers)

	v1.GET("/rules", getRules)
	v1.PATCH("/rules/:name", patchRules)
	v1.DELETE("/rules/:name", deleteRules)

	err := router.Run(":3000")
	if err != nil {
		return err
	}
	return nil
}
