package controllers

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "temp/config"
    "temp/helpers"
    "temp/models"
)

func getUserCollection() *mongo.Collection {
    return config.DB.Database("temp").Collection("users")
}

// Login handles user authentication
func Login(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    fmt.Println("Login attempt:", req.Email)
    fmt.Println("Password entered (raw):", req.Password)

    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := getUserCollection()
    err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
    if err != nil {
        fmt.Println("User not found in DB")
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    fmt.Println("Password from DB (raw):", user.Password)

    match := helpers.CheckPasswordHash(req.Password, user.Password)
    fmt.Println("Match result:", match)

    if !match {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    token, err := helpers.GenerateToken(user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}

// Signup handles user registration
func Signup(c *gin.Context) {
    var user models.User

    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    fmt.Println("Password before hashing:", user.Password)

    hashedPassword, err := helpers.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }

    fmt.Println("Hashed password to save:", hashedPassword)

    user.Password = hashedPassword
    user.ID = primitive.NewObjectID()
    user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := getUserCollection()
    _, err = collection.InsertOne(ctx, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}


// controllers/usercontroller.go
func MyProfile(c *gin.Context) {
    email := c.MustGet("email").(string)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User
    collection := getUserCollection()
    err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
        return
    }

    user.Password = "" // Hide password in response
    c.JSON(http.StatusOK, user)
}

func GetMyProfile(c *gin.Context) {
	emailInterface, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email, ok := emailInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token data"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	collection := getUserCollection()
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":           user.Name,
		"email":          user.Email,
		"skillsHave":     user.SkillsHave,
		"skillsWant":     user.SkillsWant,
		"availableDays":  user.AvailableDays,
		"createdAt":      user.CreatedAt,
	})
}
func UpdateMyProfile(c *gin.Context) {
	emailFromToken := c.MustGet("email").(string)

	var updateData struct {
		SkillsHave    []string `json:"skillsHave"`
		SkillsWant    []string `json:"skillsWant"`
		AvailableDays int      `json:"availableDays"`
	}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	collection := getUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": emailFromToken}
	update := bson.M{
		"$set": bson.M{
			"skillsHave":    updateData.SkillsHave,
			"skillsWant":    updateData.SkillsWant,
			"availableDays": updateData.AvailableDays,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
