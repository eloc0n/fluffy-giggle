package main

import (
	"fmt"
	"net/http"

	"go-api/api/dynamicvariables"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/process", func(c *gin.Context) {
		var payload map[string]interface{}
		// fmt.Printf("---------- Received payload: %v\n", payload)
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process the payload using your Go module logic
		result, err := processPayload(payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Printf("---------- skata ------------")
		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.Run(":8080")
}

func processPayload(payload map[string]interface{}) (interface{}, error) {

	// Process the "target" field based on its type
	target, ok := payload["target"]
	if !ok {
		return nil, fmt.Errorf("target field is missing")
	}

	// Assert the type of "dependencies"
	dependencies, ok := payload["dependencies"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dependencies is not a valid map[string]interface{}")
	}

	gdv := dynamicvariables.NewGenericDynamicVariables()

	var processedTarget interface{}
	switch v := target.(type) {
	case string:
		processedTarget = v
		return gdv.Inject(processedTarget, dependencies), nil
	case []interface{}:
		processedTarget = v
		return gdv.InjectIntoList(processedTarget.([]interface{}), dependencies), nil
	case map[string]interface{}:
		processedTarget = v
		return gdv.Inject(processedTarget, dependencies), nil
	default:
		return nil, fmt.Errorf("target is of an unsupported type")
	}

}
