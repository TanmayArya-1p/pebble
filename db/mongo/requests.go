package mongo

import (
	models "pebble/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateRequest(r *models.Request) (primitive.ObjectID, error) {
	r.ID = primitive.NewObjectID()
	res, err := requestsCollection.InsertOne(ctx, r)
	r.ID = res.InsertedID.(primitive.ObjectID)
	return r.ID, err
}

func GetRequest(id primitive.ObjectID) (*models.Request, error) {
	var r models.Request
	err := requestsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&r)
	return &r, err
}

func GetRequests() ([]models.Request, error) {
	var requests []models.Request
	cursor, err := requestsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var r models.Request
		cursor.Decode(&r)
		requests = append(requests, r)
	}
	return requests, nil
}

func DeleteRequest(req *models.Request) error {
	_, err := requestsCollection.DeleteOne(ctx, bson.M{"_id": req.ID})
	return err
}
