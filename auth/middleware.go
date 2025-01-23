package auth

import (
	"context"
	"fmt"
	"net/http"
	mng "pebble/db/mongo"
	rd "pebble/db/redis"
	"pebble/models"
)

func verifyUserViaSecret(uid string, secret string) (bool, models.User) {
	objId := mng.ObjIDfromString(uid)
	user, err := mng.GetUser(objId)
	if err != nil {
		fmt.Println("Could not get user with ID", uid)
		return false, models.User{}
	}
	if user.ClientSecret != secret {
		fmt.Println("Client Secret does not match")
		return false, models.User{}
	}
	rd.SetIsOnline(*user)
	return true, *user
}

func UserSecretAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("uid")
		clientSecret := r.Header.Get("secret")
		if userID == "" || clientSecret == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		stat, us := verifyUserViaSecret(userID, clientSecret)
		if !stat {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", us)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SessionCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.FormValue("sid")
		uid := r.Header.Get("uid")
		if sessionID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		objId := mng.ObjIDfromString(sessionID)
		ses, err := mng.GetSession(objId)
		userID := mng.ObjIDfromString(uid)
		user, _ := mng.GetUser(userID)
		for _, usrs := range ses.Users {
			if usrs.ID == user.ID {
				next.ServeHTTP(w, r)
				return
			}
		}

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "session", ses)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
