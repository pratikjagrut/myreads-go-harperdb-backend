package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-go-backend/database"
	"github.com/pratikjagrut/myreads-go-backend/models"
)

type BookStatus string

var (
	Reading    BookStatus = "reading"
	Finished   BookStatus = "finished"
	WantToRead BookStatus = "readlist"
)

func BookEntry(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Println("ERROR: BookEntry: ", err)
		return err
	}

	issuer, err := getIssuer(c)
	if err != nil {
		return err
	}

	book := &models.Books{
		Name:   data["name"],
		Userid: issuer,
		Status: data["status"],
	}

	tableName := database.GlobalClient.Table["books"]

	sql := fmt.Sprintf("INSERT INTO %s (name, status, userid) VALUES('%s', '%s', '%s')",
		tableName, book.Name, book.Status, book.Userid)

	log.Println("QUERY EXEC: INSERT BOOK: ", sql)
	res, err := database.GlobalClient.DB.SQLExec(sql)

	if err != nil {
		log.Println("ERROR: BookEntry: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: AddBook: %v", err),
		})
	}
	log.Println("INSERT BOOK: ", res)

	return c.JSON(book)
}

func GetBoooks(c *fiber.Ctx, which *BookStatus) error {
	var res []interface{}
	var sql string
	tableName := database.GlobalClient.Table["books"]

	issuer, err := getIssuer(c)
	if err != nil {
		return err
	}

	if which == nil {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE userid = '%s'", tableName, issuer)
	} else {
		sql = fmt.Sprintf("SELECT * FROM %s WHERE userid = '%s' AND status = '%s'", tableName, issuer, *which)
	}

	log.Println("QUERY EXEC: GetBooks: ", sql)
	err = database.GlobalClient.DB.SQLSelect(&res, sql)
	if err != nil {
		log.Println("ERROR: GetBooks: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	if len(res) == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Empty bookshelf",
		})
	}

	return c.JSON(res)
}

func UpdateStatus(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Println("ERROR: UpdateStatus: ", err)
		return err
	}

	issuer, err := getIssuer(c)
	if err != nil {
		return err
	}

	book := &models.Books{
		Name:   data["name"],
		Userid: issuer,
		Status: data["status"],
	}

	tableName := database.GlobalClient.Table["books"]
	sql := fmt.Sprintf("UPDATE %s SET status = '%s' WHERE userid = '%s' AND name = '%s'",
		tableName, book.Status, book.Userid, book.Name)

	log.Println("QUERY EXEC: UpdateStatus: ", sql)
	res, err := database.GlobalClient.DB.SQLExec(sql)

	if err != nil {
		log.Println("ERROR: UpdateStatus: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: UpdateStatus: %v", err),
		})
	}
	log.Println("UPDATE BOOK: ", res)

	return c.JSON(book)
}
