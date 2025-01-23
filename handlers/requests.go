package handlers

import (
	"encoding/json"
	"net/http"
	mng "pebble/db/mongo"
	rd "pebble/db/redis"
	"pebble/models"
)

// TODO: SUSSBBY BAKA FUNCTION
func CreateRequest(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	decoder := json.NewDecoder(r.Body)
	sid := r.FormValue("sid")
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	reqId, err := mng.CreateRequest(&req)
	if err != nil {
		http.Error(w, "Failed to append request to db", http.StatusInternalServerError)
		return
	}
	ses, _ := mng.GetSession(mng.ObjIDfromString(sid))
	ses.Requests = append(ses.Requests, reqId)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

}

func GetRequests(w http.ResponseWriter, r *http.Request) {
	sid := r.FormValue("sid")
	ses, _ := mng.GetSession(mng.ObjIDfromString(sid))
	var requests []models.Request
	for _, req := range ses.Requests {
		request, _ := mng.GetRequest(req)
		requests = append(requests, *request)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)
}

// TODO; enforce http methods
func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	reqId := r.FormValue("rid")
	req, _ := mng.GetRequest(mng.ObjIDfromString(reqId))
	err := mng.DeleteRequest(req)
	if err != nil {
		http.Error(w, "Failed to delete request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func FindSeed(w http.ResponseWriter, r *http.Request) {
	pebID := r.FormValue("pebbleID")
	idObj := mng.ObjIDfromString(pebID)

	pebble, _ := mng.GetPebble(idObj)
	for i := len(pebble.Seeds) - 1; i >= 0; i-- {
		seed := pebble.Seeds[i]
		if rd.IsOnline(seed).IsOnline {
			res := map[string]string{"seedID": seed.ID.Hex()}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}
	req := map[string]string{"seedID": "NOTFOUND"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

//Request Types:
//"BEGRTC" - intiialite webrtc tunnel
// "SRSDP" - setremotesdp
// "SASDP" - set answer sdp
// ""
