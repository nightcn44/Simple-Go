package controllers

import (
	"context"
	"net/http"
	"time"

	"bn/config"
	"bn/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var post models.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		objUserID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}
		post.AuthorID = objUserID
		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()

		collection := config.GetCollection("posts")
		result, err := collection.InsertOne(context.TODO(), post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating post"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post_id": result.InsertedID})
	}
}

func GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		collection := config.GetCollection("posts")
		var posts []models.Post

		cursor, err := collection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
			return
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var post models.Post
			cursor.Decode(&post)
			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, posts)
	}
}

func GetPostByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		collection := config.GetCollection("posts")
		var post models.Post
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&post)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		objUserID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		var post models.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := config.GetCollection("posts")
		// Check if the post exists and belongs to the current user
		var existingPost models.Post
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "author_id": objUserID}).Decode(&existingPost)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post or post not found"})
			return
		}

		update := bson.M{
			"$set": bson.M{
				"title":     post.Title,
				"content":   post.Content,
				"updated_at": time.Now(),
			},
		}

		_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			return
		}
		objUserID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			return
		}

		collection := config.GetCollection("posts")
		// Check if the post exists and belongs to the current user before deleting
		var existingPost models.Post
		err = collection.FindOne(context.TODO(), bson.M{"_id": objID, "author_id": objUserID}).Decode(&existingPost)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post or post not found"})
			return
		}

		_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}