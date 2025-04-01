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

	// Get sender's email from token
	fromEmail := c.MustGet("email").(string)
	request.ID = primitive.NewObjectID()
	request.FromEmail = fromEmail
	request.Status = "pending"
	request.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := getRequestCollection().InsertOne(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Skill request sent"})
}

func RespondToSkillRequest(c *gin.Context) {
	var body struct {
		RequestID string `json:"requestId"`
		Status    string `json:"status"` // "accepted" or "rejected"
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	requestID, err := primitive.ObjectIDFromHex(body.RequestID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = getRequestCollection().UpdateOne(ctx, bson.M{"_id": requestID}, bson.M{
		"$set": bson.M{"status": body.Status},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated"})
}
