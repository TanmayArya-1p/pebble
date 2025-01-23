package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pebble struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Seeds   []User             `bson:"seeds" json:"seeds"`
	Hash    string             `bson:"hash" json:"hash"`
	Info    string             `bson:"info" json:"info"`
	Session primitive.ObjectID `bson:"session" json:"session"`
}
