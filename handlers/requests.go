package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	mng "pebble/db/mongo"
	rd "pebble/db/redis"
	"pebble/models"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req = models.Request{
		To:      mng.ObjIDfromString(r.FormValue("to")),
		From:    mng.ObjIDfromString(r.Header.Get("uid")),
		Code:    r.FormValue("code"),
		Content: r.FormValue("content"),
	}
	sid := r.URL.Query().Get("sid")

	reqId, err := mng.CreateRequest(&req)
	if err != nil {
		http.Error(w, "Failed to append request to db", http.StatusInternalServerError)
		return
	}
	ses := r.Context().Value(models.SessionContextKey).(models.Session)
	ses.Requests = append(ses.Requests, reqId)
	res, err := mng.UpdateSession(mng.ObjIDfromString(sid), ses)
	if err != nil {
		fmt.Println("Error Updating Session")
	}
	fmt.Println(res.ModifiedCount)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

}

func GetRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ses := r.Context().Value(models.SessionContextKey).(models.Session)
	var requests []models.Request
	for _, req := range ses.Requests {
		request, _ := mng.GetRequest(req)
		requests = append(requests, *request)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)
}

func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	reqId := r.FormValue("rid")
	req, err := mng.GetRequest(mng.ObjIDfromString(reqId))
	if err != nil {
		http.Error(w, "Failed to get request", http.StatusInternalServerError)
		return
	}
	usr := r.Context().Value(models.UserContextKey).(models.User)
	if req.From != usr.ID || req.To != usr.ID {
		http.Error(w, "You are not concerned with this request", http.StatusUnauthorized)
		return
	}
	err = mng.DeleteRequest(req)
	if err != nil {
		http.Error(w, "Failed to delete request", http.StatusInternalServerError)
		return
	}
	ses := r.Context().Value(models.SessionContextKey).(models.Session)
	for i, req := range ses.Requests {
		if req == mng.ObjIDfromString(reqId) {
			ses.Requests = append(ses.Requests[:i], ses.Requests[i+1:]...)
			break
		}
	}
	_, err = mng.UpdateSession(ses.ID, ses)
	if err != nil {
		http.Error(w, "Failed to update session", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func FindSeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pebID := r.URL.Query().Get("pid")
	idObj := mng.ObjIDfromString(pebID)

	type FindSeedResponse struct {
		Found bool `json:"found"`
		Seed  models.User
	}

	pebble, _ := mng.GetPebble(idObj)
	for i := len(pebble.Seeds) - 1; i >= 0; i-- {
		seed := pebble.Seeds[i]
		if rd.IsOnline(seed).IsOnline {
			res := FindSeedResponse{
				Found: true,
				Seed:  seed,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
			return
		}
	}
	req := FindSeedResponse{
		Found: false,
		Seed:  models.User{},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

//Request Types:
//"BEGRTC" - intiialite webrtc tunnel
// "SRSDP" - setremotesdp
// "SASDP" - set answer sdp
// ""
