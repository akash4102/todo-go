package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID       primitive.ObjectID `bson:"_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
	Done     bool               `bson:"done"`
	Type     string             `bson:"type"`
	Created  primitive.DateTime `bson:"created"`
	EffortHr int                `bson:"effortHr"`
}
