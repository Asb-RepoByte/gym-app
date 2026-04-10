package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gym-app/database"
	"gym-app/models"
)

func StartSession(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	session := models.WorkoutSession{
		UserID: user.ID,
		// DepartureTime defaults to now
	}

	if result := database.DB.Create(&session); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func UpdateSession(c *gin.Context) {
	sessionID := c.Param("id")
	var body struct {
		CheckInTime    *time.Time `json:"check_in_time"`
		CheckOutTime   *time.Time `json:"check_out_time"`
		HomecomingTime *time.Time `json:"homecoming_time"`
		OverallMood    string     `json:"overall_mood"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	var session models.WorkoutSession
	if result := database.DB.First(&session, "id = ?", sessionID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	user := c.MustGet("user").(models.User)
	if session.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if body.CheckInTime != nil { session.CheckInTime = body.CheckInTime }
	if body.CheckOutTime != nil { session.CheckOutTime = body.CheckOutTime }
	if body.HomecomingTime != nil { session.HomecomingTime = body.HomecomingTime }
	if body.OverallMood != "" { session.OverallMood = body.OverallMood }

	database.DB.Save(&session)
	c.JSON(http.StatusOK, session)
}

func GetSessions(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	var sessions []models.WorkoutSession
	database.DB.Where("user_id = ?", user.ID).Find(&sessions)
	c.JSON(http.StatusOK, sessions)
}
