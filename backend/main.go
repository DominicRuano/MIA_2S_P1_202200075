package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configura CORS para permitir solicitudes desde cualquier origen
	r.Use(cors.Default())

	r.POST("/api/submit", func(c *gin.Context) {
		var json struct {
			Text string `json:"text"`
		}

		if err := c.BindJSON(&json); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "Texto recibido",
			"text":    json.Text,
		})

		fmt.Println(json.Text)
	})

	r.Run(":8080")
}
