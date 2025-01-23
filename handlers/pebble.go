package handlers

import (
	"encoding/json"
	"net/http"
	mng "pebble/db/mongo"
	"pebble/models"
)

func CreatePebble(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := r.Context().Value(models.UserContextKey).(models.User)
	session := r.Context().Value(models.SessionContextKey).(models.Session)

	hash := r.FormValue("hash")
	info := r.FormValue("info")
	seeds := []models.User{user}
	newPbl := models.Pebble{
		Hash:    hash,
		Info:    info,
		Seeds:   seeds,
		Session: session.ID,
	}
	_, err := mng.CreatePebble(&newPbl)
	if err != nil {
		http.Error(w, "Error creating pebble", http.StatusInternalServerError)
		return
	}
	session.Pebbles = append(session.Pebbles, newPbl)
	_, errU := mng.UpdateSession(session.ID, session)
	if errU != nil {
		http.Error(w, "Error updating session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newPbl)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func GetPebble(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pebbleID := r.FormValue("pid")
	sessionID := r.FormValue("sid")
	sid := mng.ObjIDfromString(sessionID)
	pid := mng.ObjIDfromString(pebbleID)
	pebble, err := mng.GetPebble(pid)
	if err != nil {
		http.Error(w, "Error getting pebble", http.StatusInternalServerError)
		return
	}
	if sid != pebble.Session {
		http.Error(w, "Pebble does not belong to session", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pebble)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func MakeMeSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pebbleID := r.FormValue("pid")
	sessionID := r.FormValue("sid")
	sid := mng.ObjIDfromString(sessionID)
	pid := mng.ObjIDfromString(pebbleID)
	pebble, err := mng.GetPebble(pid)
	if err != nil {
		http.Error(w, "Error getting pebble", http.StatusInternalServerError)
		return
	}
	if sid != pebble.Session {
		http.Error(w, "Pebble does not belong to session", http.StatusUnauthorized)
		return
	}
	pebble.Seeds = append(pebble.Seeds, models.User{ID: mng.ObjIDfromString(r.Header.Get("uid"))})
	_, err = mng.UpdatePebble(pid, *pebble)
	if err != nil {
		http.Error(w, "Error updating pebble", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pebble)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
