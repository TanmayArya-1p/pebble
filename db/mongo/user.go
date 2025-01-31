package mongo

import (
	"fmt"
	models "pebble/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

// TODO: REFERENCE OPTMIZATION INSTEAD OF STORING ENTIRE USER IN SESSION
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
			"clientsecret": updateUser.ClientSecret,
			"pwdHash":      updateUser.PwdHash,
			"localsdp":     updateUser.LocalSDP,
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
