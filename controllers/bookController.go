package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pratikjagrut/myreads-go-backend/database"
	"github.com/pratikjagrut/myreads-go-backend/models"
)

type BookStatus string

var (
	Reading  BookStatus = "reading"
	Finished BookStatus = "finished"
	Wishlist BookStatus = "wishlist"
)

func BookEntry(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("ERROR: BookEntry: Form: ", err)
		return err
	}

	if err = c.BodyParser(form); err != nil {
		log.Println("ERROR: BookEntry: BodyParser: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
			"status":  fiber.StatusInternalServerError,
		})
	}

	book := &models.Book{
		Author:      form.Value["author"][0],
		Name:        form.Value["name"][0],
		Status:      form.Value["status"][0],
		Description: form.Value["description"][0],
	}

	issuer, err := getIssuer(c)
	if err != nil {
		log.Println("ERROR: BookEntry: getIssuer", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}
	book.Userid = issuer

	tableName := database.GlobalClient.Table["books"]

	sql := fmt.Sprintf("SELECT * FROM %s WHERE userid = '%s' AND name = '%s'", tableName, book.Userid, book.Name)

	log.Println("QUERY EXEC: BookEntry: ", sql)
	var res1 []interface{}
	err = database.GlobalClient.DB.SQLSelect(&res1, sql)

	if err != nil {
		log.Println("ERROR: BookEntry: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": err,
			"status":  fiber.StatusExpectationFailed,
		})
	}

	if len(res1) != 0 {
		log.Println("ERROR: BookEntry: This book is present in your bookshelf")
		c.Status(fiber.StatusAlreadyReported)
		return c.JSON(fiber.Map{
			"message": "This book is present in your bookshelf",
			"status":  fiber.StatusAlreadyReported,
		})
	}

	new_uuid := uuid.New().String()
	images := form.File["image"]
	book.Image = fmt.Sprintf("%s_%s", new_uuid, images[0].Filename)
	staticLoc := fmt.Sprintf("./images/%s", book.Image)
	err = c.SaveFile(images[0], staticLoc)

	if err != nil {
		log.Println(err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
			"status":  fiber.StatusInternalServerError,
		})
	}

	log.Println(book)
	sql = fmt.Sprintf("INSERT INTO %s (name, status, userid, image, author, description) VALUES('%s', '%s', '%s', '%s', '%s', '%s')",
		tableName, book.Name, book.Status, book.Userid, book.Image, book.Author, book.Description)

	log.Println("QUERY EXEC: INSERT BOOK: ", sql)
	res, err := database.GlobalClient.DB.SQLExec(sql)

	if err != nil {
		log.Println("ERROR: BookEntry: SQLExec: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: AddBook: %v", err),
		})
	}
	log.Println("INSERT BOOK: ", res)
	log.Println("INSERT BOOK: Insertion successful")

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book \"%s\" added to \"%s\" section of bookshelf", book.Name, book.Status),
		"status":  fiber.StatusOK,
	})
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
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal Server Error",
			"status":  fiber.StatusInternalServerError,
		})
	}

	if len(res) == 0 {
		log.Println("ERROR: GetBooks: Empty bookshelf")
		return c.Status(fiber.StatusNotFound).SendString("Empty bookshelf")
	}

	return c.JSON(res)
}

func UpdateStatus(c *fiber.Ctx) error {
	book := new(models.Book)

	if err := c.BodyParser(book); err != nil {
		log.Println("ERROR: UpdateStatus: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	_, err := getIssuer(c)
	if err != nil {
		log.Println("ERROR: UpdateStatus: getIssuer", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	tableName := database.GlobalClient.Table["books"]
	sql := fmt.Sprintf("UPDATE %s SET status = '%s' WHERE id = '%s'",
		tableName, book.Status, book.Id)

	log.Println("QUERY EXEC: UpdateStatus: ", sql)
	res, err := database.GlobalClient.DB.SQLExec(sql)

	if err != nil {
		log.Println("ERROR: UpdateStatus: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if res.Message == "updated 0 of 0 records" {
		log.Println("ERROR: DeleteBook: ", res.Message)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "This book is not present in your bookshelf",
		})
	}

	log.Println("UPDATE BOOK: ", res)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book moved to \"%s\"", book.Status),
		"status":  fiber.StatusOK,
	})
}

func DeleteBook(c *fiber.Ctx) error {
	book := new(models.Book)

	if err := c.BodyParser(book); err != nil {
		log.Println("ERROR: DeleteBook: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	_, err := getIssuer(c)
	if err != nil {
		log.Println("ERROR: DeleteBook: getIssuer", err)
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized",
			"status":  fiber.StatusUnauthorized,
		})
	}

	tableName := database.GlobalClient.Table["books"]

	imagePath := fmt.Sprintf("./images/%s", book.Image)

	sql := fmt.Sprintf("DELETE FROM %s WHERE id = '%s'", tableName, book.Id)

	log.Println("QUERY EXEC: DeleteBook: ", sql)
	res, err := database.GlobalClient.DB.SQLExec(sql)

	if err != nil {
		log.Println("ERROR: DeleteBook: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: DeleteBook: %v", err),
		})
	}

	if res.Message == "0 of 0 records successfully deleted" {
		log.Println("ERROR: DeleteBook: ", res.Message)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "This book is not present in your bookshelf",
		})
	}

	if err = os.Remove(imagePath); err != nil {
		return err
	}

	log.Println("Delete Book: ", res)

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Book \"%s\" removed from your bookshelf.", book.Name),
		"status":  fiber.StatusOK,
	})
}
