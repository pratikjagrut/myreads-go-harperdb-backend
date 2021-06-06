package database

import (
	"log"

	"github.com/HarperDB/harperdb-sdk-go"
)

type Client struct {
	DB     *harperdb.Client
	Schema string
}

var GlobalClient *Client

func Init(host, username, password, schema string) {
	GlobalClient = &Client{
		DB:     harperdb.NewClient(host, username, password),
		Schema: schema,
	}

	err := GlobalClient.DB.CreateTable(GlobalClient.Schema, "users", "id")
	if err != nil {
		log.Println(err)
	}

	err = GlobalClient.DB.CreateTable(GlobalClient.Schema, "books", "id")
	if err != nil {
		log.Println(err)
	}
}
