package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/abe21412/slack-clone-backend/src/providers"
	"github.com/gorilla/mux"
)

type userCreatedResponse struct {
	UserID string `json:"userID"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUserID, err := providers.CreateUser(user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	log.Println(newUserID)
	json.NewEncoder(w).Encode(userCreatedResponse{newUserID})
}

func InitRouter(usersRouter *mux.Router) *mux.Router {
	usersRouter.Methods("POST").HandlerFunc(createUser)
	return usersRouter
}
