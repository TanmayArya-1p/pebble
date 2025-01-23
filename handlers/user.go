package handlers

import (
	"encoding/json"
	"net/http"
	auth "pebble/auth"
	mng "pebble/db/mongo"
	models "pebble/models"

	"github.com/google/uuid"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	newSecret := uuid.New().String()
	newUser := models.User{
		Username:     username,
		ClientSecret: newSecret,
		PwdHash:      models.GeneratePasswordHash(password),
	}
	_, err := mng.CreateUser(&newUser)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := map[string]string{"ClientSecret": newSecret, "UID": newUser.ID.Hex()}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := r.FormValue("uid")
	password := r.FormValue("password")
	if !auth.LoginVerify(password, uid) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	secy, err := auth.GetClientSecret(uid)
	if err != nil {
		http.Error(w, "Error getting client secret", http.StatusInternalServerError)
	}
	res := map[string]string{"ClientSecret": secy, "UID": uid}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
