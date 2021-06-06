package controllers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/database"
	"github.com/pratikjagrut/go-jwt-auth-harperDB/models"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := &models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	tableName := database.GlobalClient.Schema + ".users"
	sqlInsert := fmt.Sprintf("INSERT INTO %s (name, email, password) VALUES('%s', '%s', '%s')",
		tableName, user.Name, user.Email, user.Password)

	_, err := database.GlobalClient.DB.SQLExec(sqlInsert)
	if err != nil {
		err1 := fmt.Sprintf("ERROR: Register: %v", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": err1,
		})
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	tableName := database.GlobalClient.Schema + ".users"
	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", tableName, data["email"])
	var res []interface{}

	err := database.GlobalClient.DB.SQLSelect(&res, sqlStmt)

	if err != nil {
		err1 := fmt.Sprintf("ERROR: Login: %v", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": err1,
		})
	}

	if len(res) == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	var fetchedUser map[string]interface{}
	fetchedUser = res[0].(map[string]interface{})

	var user = &models.User{
		Id:       fmt.Sprintf("%v", fetchedUser["id"]),
		Email:    fmt.Sprintf("%v", fetchedUser["email"]),
		Name:     fmt.Sprintf("%v", fetchedUser["name"]),
		Password: []byte(fmt.Sprintf("%v", fetchedUser["password"])),
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect Password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login.",
		})
	}

	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Success.",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	tableName := database.GlobalClient.Schema + ".users"
	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tableName, claims.Issuer)
	var res []interface{}

	err = database.GlobalClient.DB.SQLSelect(&res, sqlStmt)

	if err != nil {
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": err,
		})
	}

	if len(res) == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User Not Found",
		})
	}

	var fetchedUser map[string]interface{}
	fetchedUser = res[0].(map[string]interface{})

	var user = &models.User{
		Id:       fmt.Sprintf("%v", fetchedUser["id"]),
		Email:    fmt.Sprintf("%v", fetchedUser["email"]),
		Name:     fmt.Sprintf("%v", fetchedUser["name"]),
		Password: []byte(fmt.Sprintf("%v", fetchedUser["password"])),
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(cookie)

	return c.JSON(fiber.Map{
		"message": "Success.",
	})
}
