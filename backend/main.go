package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/api/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello from Go!",
        })
    })
    r.Run(":8080") // Ejecuta el servidor en el puerto 8080
}
