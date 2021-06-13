package models

type Book struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Userid string `json:"userid"`
	Image  string `json:"image"`
}
