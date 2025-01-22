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
