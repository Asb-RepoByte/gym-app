package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gym-app/database"
	"gym-app/models"
)

func GetExercises(c *gin.Context) {
	var exercises []models.Exercise
	database.DB.Find(&exercises)
	c.JSON(http.StatusOK, exercises)
}

func CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&exercise); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func LogExercise(c *gin.Context) {
	var logEntry models.ExerciseLog
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&logEntry); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log exercise"})
		return
	}

	c.JSON(http.StatusOK, logEntry)
}

func AddSet(c *gin.Context) {
	var set models.Set
	if err := c.ShouldBindJSON(&set); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := database.DB.Create(&set); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add set"})
		return
	}

	c.JSON(http.StatusOK, set)
}
