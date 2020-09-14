package models

type Channel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Public      bool   `json:"public"`
	DM          bool   `json:"dm"`
	Description string `json:"description"`
}
