package main

import (
	"fmt"
	"net/http"
	"os"
	auth "pebble/auth"
	mng "pebble/db/mongo"
	rd "pebble/db/redis"
	handlers "pebble/handlers"

	"github.com/joho/godotenv"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "pong")
}

func main() {
	godotenv.Load()
	rd.Connect()
	mng.Connect()
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/user/create", handlers.CreateUser)
	http.HandleFunc("/user/login", handlers.Login)

	http.Handle("/session/join", auth.UserSecretAuth(http.HandlerFunc(handlers.JoinSession)))
	http.Handle("/session/leave", auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.LeaveSession))))

	http.Handle("/session", methodHandler{
		"POST": auth.UserSecretAuth(http.HandlerFunc(handlers.CreateSession)),
		"GET":  auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.SessionMetadata))),
	})

	http.Handle("/request", methodHandler{
		"POST":   auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.CreateRequest))),
		"DELETE": auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.DeleteRequest))),
		"GET":    auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.GetRequests))),
	})

	http.Handle("/pebble/findSeed", auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.FindSeed))))
	http.Handle("/pebble", methodHandler{
		"POST": auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.CreatePebble))),
		"GET":  auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.GetPebble))),
	})

	http.Handle("/pebble/mms", auth.UserSecretAuth(auth.SessionCheck(http.HandlerFunc(handlers.MakeMeSeed))))

	fmt.Println("Starting server at port " + os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

type methodHandler map[string]http.Handler

func (h methodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := h[r.Method]; ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
