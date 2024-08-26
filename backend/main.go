package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//	Configura CORS para permitir solicitudes desde el origen del frontend
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Cambia esto si tu frontend está en otra dirección
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
		fmt.Println("Texto recibido!", inputText)

		// Procesamiento del texto (ejemplo: convertir a mayúsculas)
		processedText := processText(inputText)

		c.JSON(http.StatusOK, gin.H{
			"message": "Texto procesado con éxito",
			"text":    processedText,
		})
	})

	r.Run(":8080")
}

// Función de ejemplo para procesar el texto
func processText(text string) string {
	// Aquí puedes implementar la lógica que desees
	processed := strings.ToUpper(text)
	return processed
}
