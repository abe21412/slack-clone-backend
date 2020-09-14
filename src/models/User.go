package models

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}
