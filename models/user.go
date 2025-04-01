package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name          string             `bson:"name" json:"name"`
    Email         string             `bson:"email" json:"email"`
    Password string                  `bson:"password" json:"password"`
    SkillsHave    []string           `bson:"skillsHave" json:"skillsHave"`
    SkillsWant    []string           `bson:"skillsWant" json:"skillsWant"`
    AvailableDays int                `bson:"availableDays" json:"availableDays"`
    CreatedAt     primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
