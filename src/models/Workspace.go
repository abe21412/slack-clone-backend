package models

type Workspace struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner"`
}
