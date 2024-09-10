package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	analyzer "Backend/Analyzer"
)

func main() {
	r := gin.Default()

	//	Configuracion de CORS para permitir solicitudes desde el origen del frontend
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}

	r.Use(cors.New(config))

	r.POST("/api/submit", func(c *gin.Context) {
		var requestBody struct {
			Text string `json:"text"`
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de solicitud inválido"})
			return
		}

		inputText := requestBody.Text

		// Analizador	de texto
		processedText := analyzer.Analyzer(inputText)

		c.JSON(http.StatusOK, gin.H{
			"message": "Texto procesado con éxito",
			"text":    processedText,
		})
	})

	r.Run(":8080")
}
