package controllers

import (
	"crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInput struct {
	Email  string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

type UserInputResponse struct {
	Email string `json:"email"`
	ID uuid.UUID `json:"id"`
}

func Register(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New()
	user := models.User{ID: id, Email: input.Email, Password: input.Password}
	db := c.MustGet("db").(*gorm.DB)
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exist"})
		return
	}
	response := UserInputResponse{ID: id, Email: input.Email}
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	db.Find(&user)
	response := UserInputResponse{ID: user.ID, Email: user.Email}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func GetUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)
	response := []UserInputResponse{}
	for _, user := range users {
		response = append(response, UserInputResponse{user.Email, user.ID})
	}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	update := models.User{ID: user.ID, Email: input.Email, Password: input.Password}
	db.Save(&update)
	c.JSON(http.StatusOK, gin.H{"data": "Record updated"})
}
