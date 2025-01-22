package mongo

import (
	"fmt"
	models "pebble/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePebble(pebble *models.Pebble) (primitive.ObjectID, error) {
	pebble.ID = primitive.NewObjectID()
	result, err := pebblesCollection.InsertOne(ctx, *pebble)
	if err != nil {
		fmt.Println("Error creating pebble:", err)
		return primitive.NewObjectID(), err
	}
	pebble.ID = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetPebble(id primitive.ObjectID) (*models.Pebble, error) {
	var pebble models.Pebble
	fmt.Println("Searching for pebble with ID:", id)
	err := pebblesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&pebble)
	if err != nil {
		return nil, err
	}
	return &pebble, nil
}

func UpdatePebble(id primitive.ObjectID, updatePebble models.Pebble) (*mongo.UpdateResult, error) {
	update := bson.M{
		"$set": bson.M{
			"seeds": updatePebble.Seeds,
			"hash":  updatePebble.Hash,
			"info":  updatePebble.Info,
		},
	}
	result, err := pebblesCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		fmt.Println("Error Updating Pebble", err)
		return nil, err
	}
	return result, nil
}

func DeletePebble(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := pebblesCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		fmt.Println("Error Deleting Pebble", err)
		return nil, err
	}
	return result, nil
}
