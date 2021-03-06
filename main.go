package main

import (
	"flatman-api/controllers"
	"flatman-api/middlewares"
	"flatman-api/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDataBase()

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	public := r.Group("/api")

	public.POST("/login", controllers.Login)

	public.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": "1.0"})
	})

	private := r.Group("/api")

	private.Use(middlewares.JwtAuthMiddleware())
	{
		// User Info
		private.GET("/user", controllers.CurrentUser)

		// Flat API
		private.GET("/flat", controllers.GetFlats)
		private.GET("/flat/:id", controllers.GetFlat)
		private.POST("/flat", controllers.SaveFlat)
		private.PUT("/flat/:id", controllers.UpdateFlat)
		private.DELETE("/flat/:id", controllers.DeleteFlat)

		// Lessee API
		private.GET("/lessee", controllers.GetLessees)
		private.GET("/lessee/:id", controllers.GetLessee)
		private.POST("/lessee", controllers.SaveLessee)
		private.PUT("/lessee/:id", controllers.UpdateLessee)
		private.DELETE("/lessee/:id", controllers.DeleteLessee)

		// Balance API
		private.GET("/balance", controllers.GetBalances)
		private.GET("/balance/:id", controllers.GetBalance)
		private.POST("/balance", controllers.SaveBalance)
		private.PUT("/balance/:id", controllers.UpdateBalance)
		private.DELETE("/balance/:id", controllers.DeleteBalance)
	}

	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalf("Error starting server")
	}
}
