package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	ID       primitive.ObjectID   `bson:"_id" json:"id"`
	Key      string               `json:"key"`
	Pebbles  []Pebble             `bson:"pebbles" json:"pebbles"`
	Users    []User               `bson:"users" json:"users"`
	Requests []primitive.ObjectID `bson:"requests" json:"requests"`
}
