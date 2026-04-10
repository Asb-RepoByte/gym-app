package main

import (
	"github.com/gin-gonic/gin"
	"gym-app/controllers"
	"gym-app/database"
	"gym-app/middleware"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	api := r.Group("/api")
	api.Use(middleware.RequireAuth)
	{
		api.GET("/me", func(c *gin.Context) {
			user := c.MustGet("user")
			c.JSON(200, user)
		})
		
		api.POST("/sessions", controllers.StartSession)
		api.PUT("/sessions/:id", controllers.UpdateSession)
		api.GET("/sessions", controllers.GetSessions)

		api.GET("/exercises", controllers.GetExercises)
		api.POST("/exercises", controllers.CreateExercise)

		api.POST("/exercises/log", controllers.LogExercise)
		api.POST("/exercises/sets", controllers.AddSet)
	}

	r.Run(":8000")
}
