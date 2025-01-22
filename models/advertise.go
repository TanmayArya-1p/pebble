package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	To      primitive.ObjectID `bson:"to" json:"to"`
	From    primitive.ObjectID `bson:"from" json:"from"`
	Content string             `bson:"content" json:"content"`
}
