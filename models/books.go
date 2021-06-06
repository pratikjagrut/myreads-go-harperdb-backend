package models

type Books struct {
	Id     string `json:id`
	Name   string `json:name`
	Status string `json:status`
	Userid string `json:userid`
}
