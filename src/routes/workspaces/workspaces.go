package workspaces

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/abe21412/slack-clone-backend/src/models"
	"github.com/abe21412/slack-clone-backend/src/providers"
	"github.com/gorilla/mux"
)

type workspaceCreatedResponse struct {
	WorkspaceID string `json:"workspaceID"`
}

type workspacesGetResponse struct {
	Workspaces []models.Workspace `json:"workspaces"`
}

func createWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace *models.Workspace
	err := json.NewDecoder(r.Body).Decode(&workspace)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newWorkspaceID, err := providers.CreateWorkspace(workspace)
	if err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "duplicate") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println(newWorkspaceID)
	json.NewEncoder(w).Encode(workspaceCreatedResponse{newWorkspaceID})
}

func getWorkspaces(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "no userID passed in", http.StatusBadRequest)
		return
	}
	workspaces, err := providers.GetWorkspaces(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(workspaces)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workspacesGetResponse{workspaces})
}

func InitRouter(workspacesRouter *mux.Router) *mux.Router {
	workspacesRouter.Methods("POST").HandlerFunc(createWorkspace)
	workspacesRouter.Methods("GET").Queries("userID", "{userID}").HandlerFunc(getWorkspaces)
	return workspacesRouter
}
