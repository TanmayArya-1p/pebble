package main

import (
	"context"
	"fmt"
	"os"
	models "pebble/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var pebblesCollection *mongo.Collection
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
	sessionsCollection = client.Database("pebble").Collection("sessions")
	pebblesCollection = client.Database("pebble").Collection("pebbles")

}

func CreateUser(user *models.User) (primitive.ObjectID, error) {
	user.ID = primitive.NewObjectID()
	result, err := usersCollection.InsertOne(ctx, *user)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return primitive.NewObjectID(), err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
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

func main() {
	godotenv.Load()
	Connect()
	testuser := models.User{
		Username:     "testuser",
		ClientSecret: "da",
		PwdHash:      "da",
	}
	res, _ := CreateUser(&testuser)

	tus, _ := GetUser(res)
	tsession := models.Session{
		Key:     "123",
		Pebbles: []models.Pebble{},
		Users:   []models.User{*tus},
	}
	CreateSession(&tsession)
}
