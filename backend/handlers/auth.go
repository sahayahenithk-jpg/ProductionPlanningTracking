package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"projectplanningtracking/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username"`
}

type loginInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input registerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		input.Email = strings.ToLower(strings.TrimSpace(input.Email))
		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" && strings.TrimSpace(input.Username) != "" {
			input.Name = strings.TrimSpace(input.Username)
		}

		var existing models.User
		if err := db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to hash password"})
			return
		}

		user := models.User{
			Name:         input.Name,
			Email:        input.Email,
			PasswordHash: string(hashedPassword),
			Role:         "operator",
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input loginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		identifier := strings.TrimSpace(input.Username)
		if identifier == "" {
			identifier = strings.TrimSpace(input.Email)
		}
		if identifier == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username or email is required"})
			return
		}

		var user models.User
		if err := db.Where("LOWER(email) = LOWER(?) OR LOWER(name) = LOWER(?)", identifier, identifier).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if strings.TrimSpace(user.Role) == "" {
			user.Role = "operator"
			_ = db.Model(&user).Update("role", user.Role).Error
		}

		secret := getEnv("JWT_SECRET", "supersecretkey")
		claims := JwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   strconv.Itoa(int(user.ID)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Role: strings.ToLower(user.Role),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": signedToken, "role": user.Role})
	}
}

func Profile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, err := strconv.ParseUint(userIDRaw.(string), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		role := strings.ToLower(user.Role)
		if role == "" {
			role = "operator"
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  role,
		})
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
