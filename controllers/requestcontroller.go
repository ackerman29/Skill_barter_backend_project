
package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"temp/config"
	"temp/models"
)

func getRequestCollection() *mongo.Collection {
	return config.DB.Database("temp").Collection("requests")
}

func SendSkillRequest(c *gin.Context) {
	var request models.SkillRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	fromEmail := c.MustGet("email").(string)

	// Get sender's name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var sender models.User
	err := config.DB.Database("temp").Collection("users").FindOne(ctx, bson.M{"email": fromEmail}).Decode(&sender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sender not found"})
		return
	}

	request.ID = primitive.NewObjectID()
	request.FromEmail = fromEmail
	request.FromName = sender.Name
	request.Status = "pending"
	request.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = getRequestCollection().InsertOne(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Skill request sent"})
}


func RespondToSkillRequest(c *gin.Context) {
	var body struct {
		FromName string `json:"fromName"`
		Status   string `json:"status"` // "accepted" or"rejected"
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	toEmail := c.MustGet("email").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"fromName": body.FromName,
		"toEmail":  toEmail,
		"status":   "pending",
	}

	update := bson.M{"$set": bson.M{"status": body.Status}}

	result, err := getRequestCollection().UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request (maybe not found or already handled)"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated"})
}

