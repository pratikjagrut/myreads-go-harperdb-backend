package database

import (
	"fmt"
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

func (c *Client) Where(table, attribute, value string) ([]interface{}, error) {
	var res []interface{}

	tableName := fmt.Sprintf("%v.%v", c.Schema, table)
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", tableName, attribute, value)

	log.Println("Query executing: ", sql)

	err := c.DB.SQLSelect(&res, sql)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Insert(table string, m map[string]string) (*harperdb.AffectedResponse, error) {
	attribute := make([]string, 0, len(m))
	values := make([]string, 0, len(m))

	for k, v := range m {
		attribute = append(attribute, k)
		values = append(values, v)
	}

	tableName := fmt.Sprintf("%v.%v", c.Schema, table)
	sqlInsert := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES('%s', '%s', '%s')",
		tableName, attribute[0], attribute[1], attribute[2], values[0], values[1], values[2])

	res, err := c.DB.SQLExec(sqlInsert)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectAll(table string) ([]interface{}, error) {
	var res []interface{}

	tableName := fmt.Sprintf("%v.%v", c.Schema, table)
	sql := fmt.Sprintf("SELECT * FROM %s", tableName)

	log.Println("Query executing: ", sql)

	err := c.DB.SQLSelect(&res, sql)
	if err != nil {
		return nil, err
	}

	return res, nil
}
