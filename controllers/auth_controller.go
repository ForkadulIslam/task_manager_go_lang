package controllers

import (
	"bytes"
	"encoding/json"
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

type ExternalLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
	Name    string `json:"name"`
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

	// Prepare the request to the external API
	requestBody, err := json.Marshal(map[string]string{
		"username": input.Username,
		"password": input.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
		return
	}

	resp, err := http.Post("https://member.techvengersltd.com/api/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external login API"})
		return
	}
	defer resp.Body.Close()

	var externalResponse ExternalLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&externalResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode external login response"})
		return
	}

	if !externalResponse.Success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": externalResponse.Message})
		return
	}

	// If login is successful, find or create the user in the local DB
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// User not found, create a new one
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password for new user"})
			return
		}
		newUser := models.User{
			Username: input.Username,
			Password: string(hashedPassword),
			// Status and UserLabel are not in the login response, so I'll use defaults.
			Status:    1, // Active
			UserLabel: 2, // User
		}
		if err := database.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new user"})
			return
		}
		user = newUser
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"user_id":    user.ID,
		"user_label": user.UserLabel,
		"username":   user.Username,
	})
}

type RemoteUser struct {
	Username  string `json:"username"`
	Status    int    `json:"status"`
	UserLabel int    `json:"user_label"`
}

func SyscUser(c *gin.Context) {
	resp, err := http.Get("https://member.techvengersltd.com/api/get-users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users from the external API"})
		return
	}
	defer resp.Body.Close()

	var remoteUsers []RemoteUser
	if err := json.NewDecoder(resp.Body).Decode(&remoteUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users from the external API"})
		return
	}

	for _, remoteUser := range remoteUsers {
		var user models.User
		if err := database.DB.Where("username = ?", remoteUser.Username).First(&user).Error; err != nil {
			// User not found, create a new one
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			if err != nil {
				// Log the error, but continue with other users
				continue
			}
			newUser := models.User{
				Username:  remoteUser.Username,
				Password:  string(hashedPassword),
				Status:    remoteUser.Status,
				UserLabel: remoteUser.UserLabel,
			}
			database.DB.Create(&newUser)
		} else {
			// User found, update status and user_label
			user.Status = remoteUser.Status
			user.UserLabel = remoteUser.UserLabel
			database.DB.Save(&user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User synchronization completed successfully"})
}
