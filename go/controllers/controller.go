package controllers

import (
	"jwt/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

func Register(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	id := uuid.New()
	user := models.User{ID: id, Email: input.Email, Password: input.Password}
	db := c.MustGet("db").(*gorm.DB)
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exist"})
		return
	}
	response := UserInputResponse{ID: id, Email: input.Email}
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

func Login(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?",input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User doesn't exist"})
		return
	}
	if input.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong password"})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: input.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	token := c.Request.Header.Get("x-auth-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "JWT is missing"})
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	result := db.Create(&models.Token{Token: token})
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Token already exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func Welcome(c *gin.Context) {
	token := c.Request.Header.Get("x-auth-token")
	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "JWT is missing"})
		return
	}
	var revokeToken models.Token
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("token = ?",token).First(&revokeToken).Error; err == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "The token has been revoked"})
		return
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if time.Second > time.Until(claims.ExpiresAt.Time){
			c.JSON(http.StatusUnauthorized, gin.H{"message": "JWT is expired"})
			return
		}
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "JWT signature is invalid"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid JWT token"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Welcome "+claims.Username})
}