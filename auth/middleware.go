package auth

import (
	"fmt"
	"net/http"
	mng "pebble/db/mongo"
)

func verifyUserViaSecret(uid string, secret string) bool {
	objId := mng.ObjIDfromString(uid)
	user, err := mng.GetUser(objId)
	if err != nil {
		fmt.Println("Could not get user with ID", uid)
		return false
	}
	if user.ClientSecret != secret {
		fmt.Println("Client Secret does not match")
		return false
	}
	return true
}

func UserSecretAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("uid")
		clientSecret := r.Header.Get("secret")
		if userID == "" || clientSecret == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !verifyUserViaSecret(userID, clientSecret) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
