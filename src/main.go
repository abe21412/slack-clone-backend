package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abe21412/slack-clone-backend/src/routes/channels"

	"github.com/abe21412/slack-clone-backend/src/routes/workspaces"

	"github.com/abe21412/slack-clone-backend/src/routes/users"
	"github.com/abe21412/slack-clone-backend/src/util/db"
	"github.com/abe21412/slack-clone-backend/src/ws"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db.Init()
	defer db.Close()
	hub := ws.NewHub()
	go hub.Run()
	//cors := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	//cors2 := handlers.AllowedHeaders([]string{"*"})
	//cors3 := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("[1,2,3]")
	})
	api := r.PathPrefix("/api").Subrouter()

	usersRouter := api.PathPrefix("/users").Subrouter()
	users.InitRouter(usersRouter)

	workspacesRouter := api.PathPrefix("/workspaces").Subrouter()
	workspaces.InitRouter(workspacesRouter)

	channelsRouter := api.PathPrefix("/channels").Subrouter()
	channels.InitRouter(channelsRouter)

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
