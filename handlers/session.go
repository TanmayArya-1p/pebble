package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	mng "pebble/db/mongo"
	"pebble/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value(models.UserContextKey).(models.User)
	newSes := models.Session{
		Key:      r.FormValue("key"),
		Pebbles:  []models.Pebble{},
		Users:    []models.User{user},
		Requests: []primitive.ObjectID{},
	}
	sid, error := mng.CreateSession(&newSes)
	if error != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := map[string]string{"SessionID": sid.Hex()}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("Error encoding response", err)
	}
}

func JoinSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid := r.URL.Query().Get("sid")
	key := r.FormValue("key")
	sessionID := mng.ObjIDfromString(sid)
	session, err := mng.GetSession(sessionID)
	if err != nil {
		http.Error(w, "Error getting session", http.StatusInternalServerError)
		return
	}
	user := r.Context().Value(models.UserContextKey).(models.User)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}
	if session.Key != key {
		http.Error(w, "Invalid key", http.StatusUnauthorized)
		return
	}
	for _, u := range session.Users {
		if u.ID == user.ID {
			http.Error(w, "User already in session", http.StatusConflict)
			return
		}
	}
	session.Users = append(session.Users, user)
	_, err = mng.UpdateSession(sessionID, *session)
	if err != nil {
		http.Error(w, "Error updating session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		fmt.Println("Error encoding response", err)
	}
}

type ResponseSession struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Key      string             `json:"key"`
	Pebbles  []models.Pebble    `bson:"pebbles" json:"pebbles"`
	Users    []models.User      `bson:"users" json:"users"`
	Requests []models.Request   `bson:"requests" json:"requests"`
}

func SessionMetadata(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println(r.Context().Value(models.SessionContextKey))
	session := r.Context().Value(models.SessionContextKey).(models.Session)
	localSDP := r.FormValue("localSDP")
	usr := r.Context().Value(models.UserContextKey).(models.User)
	if localSDP != "" {
		usr.LocalSDP = localSDP
		_, err := mng.UpdateUser(usr.ID, usr)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseSession := ResponseSession{
		ID:       session.ID,
		Key:      session.Key,
		Pebbles:  session.Pebbles,
		Users:    session.Users,
		Requests: []models.Request{},
	}

	for _, reqID := range session.Requests {
		req, err := mng.GetRequest(reqID)
		if err != nil {
			http.Error(w, "Error getting request", http.StatusInternalServerError)
			return
		}
		responseSession.Requests = append(responseSession.Requests, *req)
	}

	err := json.NewEncoder(w).Encode(responseSession)
	if err != nil {
		fmt.Println("Error encoding response", err)
	}
}

func LeaveSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session := r.Context().Value(models.SessionContextKey).(models.Session)
	user := r.Context().Value(models.UserContextKey).(models.User)
	for i, us := range session.Users {
		if us.ID == user.ID {
			session.Users = append(session.Users[:i], session.Users[i+1:]...)
			break
		}
	}
	_, err := mng.UpdateSession(session.ID, session)
	if err != nil {
		http.Error(w, "Error updating session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(session)
	if err != nil {
		fmt.Println("Error encoding response", err)
	}
}
