package main

import (
	"fmt"
	"net/http"
	"strings"

	"go-api/api/dynamicvariables"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is missing"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Token parsing error: %v", err)})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	secretKey := "b7ac1ebe1d3b143385067230b4db785cbdc4a4cfc1603ee548c79f69a94c00ef"

	authorized := router.Group("/")
	authorized.Use(AuthMiddleware(secretKey))
	{
		authorized.POST("/process", processHandler)
	}

	router.Run(":8080")
}

func processHandler(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := processPayload(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func processPayload(payload map[string]interface{}) (interface{}, error) {
	target, ok := payload["target"]
	if !ok {
		return nil, fmt.Errorf("target field is missing")
	}

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
