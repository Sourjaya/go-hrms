package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"` //in mongo db ID is stored with an underscore at front
	Name   string             `json:"name" bson:"name"`
	Salary string             `json:"salary" bson:"salary"`
	Age    int                `json:"age" bson:"age"`
}
