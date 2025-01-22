package mongo

import (
	"fmt"
	models "pebble/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateSession(session *models.Session) (primitive.ObjectID, error) {
	session.ID = primitive.NewObjectID()
	result, err := sessionsCollection.InsertOne(ctx, *session)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return primitive.NewObjectID(), err
	}
	session.ID = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetSession(id primitive.ObjectID) (*models.Session, error) {
	var session models.Session
	fmt.Println("Searching for session with ID:", id)
	err := sessionsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func UpdateSession(id primitive.ObjectID, updateSession models.Session) (*mongo.UpdateResult, error) {
	update := bson.M{
		"$set": bson.M{
			"key":     updateSession.Key,
			"pebbles": updateSession.Pebbles,
			"users":   updateSession.Users,
		},
	}
	result, err := sessionsCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		fmt.Println("Error Updating Session", err)
		return nil, err
	}
	return result, nil
}

func DeleteSession(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := sessionsCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		fmt.Println("Error Deleting Session", err)
		return nil, err
	}
	return result, nil
}
