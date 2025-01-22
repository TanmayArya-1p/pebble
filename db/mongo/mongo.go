package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var pebblesCollection *mongo.Collection
var requestsCollection *mongo.Collection
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
	requestsCollection = client.Database("pebble").Collection("requests")
}

func ObjIDfromString(id string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error converting string to ObjectID", err)
	}
	return objID
}

// func main() {
// 	godotenv.Load()
// 	Connect()
// 	testuser := models.User{
// 		Username:     "testuser",
// 		ClientSecret: "da",
// 		PwdHash:      "da",
// 	}
// 	res, _ := CreateUser(&testuser)

// 	tus, _ := GetUser(res)
// 	tsession := models.Session{
// 		Key:     "123",
// 		Pebbles: []models.Pebble{},
// 		Users:   []models.User{*tus},
// 	}
// 	CreateSession(&tsession)
// }
