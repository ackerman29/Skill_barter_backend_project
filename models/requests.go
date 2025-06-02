package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SkillRequest struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FromEmail  string             `bson:"fromEmail" json:"fromEmail"`

	FromName  string             `bson:"fromName" json:"fromName"`

	ToEmail    string             `bson:"toEmail" json:"toEmail"`
	Skill      string             `bson:"skill" json:"skill"`
	Status     string             `bson:"status" json:"status"` // pending, accepted, rejected
	CreatedAt  primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
