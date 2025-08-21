package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-sql-driver/mysql"

	"taskmanager/database"
	"taskmanager/models"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			_ = ve.Translate(trans)

			for _, e := range ve {
				if e.Field() == "Username" && e.Tag() == "required" {
					errors = append(errors, "Username is required")
				} else if e.Field() == "Password" && e.Tag() == "required" {
					errors = append(errors, "Password is required")
				} else {
					errors = append(errors, e.Translate(trans))
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to hash password"}})
		return	}

	user := models.User{
		Username:  input.Username,
		Password:  string(hashedPassword),
		UserLabel: 2, // Default to regular user
	}

	if err := database.DB.Create(&user).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusConflict, gin.H{"errors": []string{"Username already exists"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to create user"}})
		return	}


	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []string
		if ve, ok := err.(validator.ValidationErrors); ok {
			en := en.New()
			uni := ut.New(en, en)
			trans, _ := uni.GetTranslator("en")

			_ = ve.Translate(trans)

			for _, e := range ve {
				if e.Field() == "Username" && e.Tag() == "required" {
					errors = append(errors, "Username is required")
				} else if e.Field() == "Password" && e.Tag() == "required" {
					errors = append(errors, "Password is required")
				} else {
					errors = append(errors, e.Translate(trans))
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte("your_secret_key")) // Replace with a strong secret key
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
