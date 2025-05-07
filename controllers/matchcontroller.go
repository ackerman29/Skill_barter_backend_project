package controllers
import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"temp/config"
	"temp/models"
)

func MatchUsers(c *gin.Context) {
	email := c.MustGet("email").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var currentUser models.User
	var matches []models.User
	var err1, err2 error

	collection := config.DB.Database("temp").Collection("users")

	// Run user fetch in one goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		err1 = collection.FindOne(ctx, bson.M{"email": email}).Decode(&currentUser)
	}()

	// Run match fetch in another goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		cursor, err := collection.Find(ctx, bson.M{
			"skillsHave": bson.M{"$in": currentUser.SkillsWant},
			"skillsWant": bson.M{"$in": currentUser.SkillsHave},
			"email":      bson.M{"$ne": email},
		})
		if err != nil {
			err2 = err
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var user models.User
			if err := cursor.Decode(&user); err == nil {
				matches = append(matches, user)
			}
		}
	}()

	// Wait for both to complete
	wg.Wait()

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": currentUser,
		"matches": matches,
	})
}
