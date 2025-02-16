package auth

import (
	"fmt"
	mng "pebble/db/mongo"
)

func LoginVerify(passwd string, uid string) bool {
	objId := mng.ObjIDfromString(uid)
	user, error := mng.GetUser(objId)
	if error != nil {
		fmt.Println("Error getting user with ID", uid)
		return false
	}
	if !user.VerifyPassword(passwd) {
		fmt.Println("Password does not match")
		return false
	}
	return true
}

func GetClientSecret(uid string) (string, string, error) {
	objId := mng.ObjIDfromString(uid)
	user, error := mng.GetUser(objId)
	if error != nil {
		fmt.Println("Error getting user with ID", uid)
		return "", "", error
	}
	return user.ClientSecret, user.InSession, nil
}
