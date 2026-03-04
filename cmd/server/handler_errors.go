package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func handlerError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			log.Printf("Error: %v", err.Error())
		}
	}
}
