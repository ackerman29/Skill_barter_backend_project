package controllers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"temp/models"
	"temp/config"
)

func MatchUsers(c *gin.Context) {
	email := c.MustGet("email").(string)
	collection := config.DB.Database("temp").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var currentUser models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Query to find matching users
	cursor, err := collection.Find(ctx, bson.M{
		"skillsHave": bson.M{"$in": currentUser.SkillsWant},
		"skillsWant": bson.M{"$in": currentUser.SkillsHave},
		"email":      bson.M{"$ne": email},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding matches"})
		return
	}
	defer cursor.Close(ctx)

	var (
		matchNames []string
		mu         sync.Mutex
		wg         sync.WaitGroup
	)

	// Loop through cursor results in main goroutine
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		// Spawn a goroutine to process user concurrently
		wg.Add(1)
		go func(u models.User) {
			defer wg.Done()

			// Safely append to matchNames
			mu.Lock()
			matchNames = append(matchNames, u.Name)
			mu.Unlock()
		}(user)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{
		"matchCount": len(matchNames),
		"names":      matchNames,
	})
}
