package controllers

import (
	"bab7/database"
	"bab7/entities"
	"bab7/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func InitProj(c *gin.Context) {
	err := repositories.InitProj(database.DbConnection)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetStatus(c *gin.Context) {
	var result gin.H
	status, err := repositories.GetStatus(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err.Error(),
		}
	} else {
		result = gin.H{
			"result": status,
		}
	}
	c.JSON(http.StatusOK, result)
}

func UpdateStatus(c *gin.Context) {
	var status entities.Servo
	srv_status, _ := strconv.Atoi(c.Param("srv_status"))
	status.StatusServo = srv_status
	err := repositories.UpdateStatus(database.DbConnection, status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"StatusServo": status.StatusServo})
}
