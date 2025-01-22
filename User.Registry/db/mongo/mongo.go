package mongo

import (
	"context"
	"fmt"
	"os"
	models "pebble/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var ctx = context.Background()

func Connect() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_STRING"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB", err)
	}
	usersCollection = client.Database("pebble").Collection("users")
}

func CreateUser(user *models.User) (*mongo.InsertOneResult, error) {
	user.ID = primitive.NewObjectID()
	result, err := usersCollection.InsertOne(ctx, *user)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return result, nil
}

func GetUser(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	fmt.Println("Searching for user with ID:", id)
	err := usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id primitive.ObjectID, updateUser models.User) (*mongo.UpdateResult, error) {
	update := bson.M{
		"$set": bson.M{
			"username":     updateUser.Username,
			"clientSecret": updateUser.ClientSecret,
			"pwdHash":      updateUser.PwdHash,
		},
	}
	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		fmt.Println("Error Updating User", err)
		return nil, err
	}
	return result, nil
}

func DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := usersCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		fmt.Println("Error Deleting User", err)
		return nil, err
	}
	return result, nil
}

func ObjIDfromString(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error converting string to ObjectID", err)
	}
	return objID
}
