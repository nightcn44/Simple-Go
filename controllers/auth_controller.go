package controllers

import (
	"context"
	"net/http"
	"time"

	"bn/config"
	"bn/models"
	"bn/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if username or email already exists
		collection := config.GetCollection("users")
		var existingUser models.User
		err := collection.FindOne(context.TODO(), bson.M{"$or": []bson.M{{"username": user.Username}, {"email": user.Email}}}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already registered"})
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		result, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		token, err := utils.GenerateToken(result.InsertedID.(primitive.ObjectID).Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": token})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginCreds struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&loginCreds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := config.GetCollection("users")
		var user models.User
		err := collection.FindOne(context.TODO(), bson.M{"email": loginCreds.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if !utils.CheckPasswordHash(loginCreds.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, err := utils.GenerateToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
	}
}
