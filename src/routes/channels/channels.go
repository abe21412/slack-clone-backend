package channels

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/abe21412/slack-clone-backend/src/providers"
	"github.com/gorilla/mux"
)

type createChannelRequest struct {
	Channel *models.Channel `json:"channel"`
	UserID  string          `json:"userID"`
}

type channelsGetResponse struct {
	Channels []models.Channel `json:"channels"`
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	var req createChannelRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newChannel, err := providers.CreateChannel(req.Channel, req.UserID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println(newChannel)
	json.NewEncoder(w).Encode(newChannel)
}

func getChannels(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	workspaceID := r.URL.Query().Get("workspaceID")
	if userID == "" || workspaceID == "" {
		http.Error(w, "please pass in both userID and workspaceID", http.StatusBadRequest)
		return
	}
	channels, err := providers.GetChannels(userID, workspaceID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(channels)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(channelsGetResponse{channels})
}

func InitRouter(channelsRouter *mux.Router) *mux.Router {
	channelsRouter.Methods("POST").HandlerFunc(createChannel)
	channelsRouter.Methods("GET").Queries("userID", "{userID}", "workspaceID", "{workspaceID}").HandlerFunc(getChannels)
	return channelsRouter
}
