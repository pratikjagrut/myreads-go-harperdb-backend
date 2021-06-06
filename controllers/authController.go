package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/pratikjagrut/myreads-go-backend/database"
	"github.com/pratikjagrut/myreads-go-backend/models"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Println("ERROR: Register: ", err)
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := &models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	tableName := database.GlobalClient.Table["users"]
	sql := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", tableName, user.Email)

	log.Println("REGISTER: QUERY EXEC: ", sql)
	var res []interface{}
	err := database.GlobalClient.DB.SQLSelect(&res, sql)

	if err != nil {
		log.Println("ERROR: Register: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: Register: %v", err),
		})
	}

	if len(res) != 0 {
		log.Println("This email id is already registered.")
		c.Status(fiber.StatusFound)
		return c.JSON(fiber.Map{
			"message": "This email id is already registered.",
		})
	}

	sqlInsert := fmt.Sprintf("INSERT INTO %s (email, name, password) VALUES('%s', '%s', '%s')",
		tableName, user.Email, user.Name, user.Password)

	log.Println("REGISTER: QUERY EXEC: ", sqlInsert)
	res1, err := database.GlobalClient.DB.SQLExec(sqlInsert)

	if err != nil {
		log.Println("ERROR: Register: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: Register: %v", err),
		})
	}
	log.Println("REGISTER: INSERT RES: ", res1)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		log.Println("ERROR: Login: ", err)
		return err
	}

	tableName := database.GlobalClient.Table["users"]
	sql := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", tableName, data["email"])

	log.Println("LOGIN: QUERY EXEC: ", sql)
	var res []interface{}
	err := database.GlobalClient.DB.SQLSelect(&res, sql)

	if err != nil {
		log.Println("ERROR: Login: ", err)
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: Login: %v", err),
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
		log.Println("ERROR: Login: ", err)
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
		log.Println("ERROR: Login: ", err)
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
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
		"message": "Login Success",
	})
}

func getIssuer(c *fiber.Ctx) (string, error) {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		log.Println("ERROR: getIssuer: ", err)
		c.Status(fiber.StatusUnauthorized)
		return "", c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil

}

func User(c *fiber.Ctx) error {
	issuer, err := getIssuer(c)
	if err != nil {
		return err
	}

	tableName := database.GlobalClient.Table["users"]
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id = '%s'", tableName, issuer)

	log.Println("USER: QUERY EXEC: ", sql)
	var res []interface{}
	err = database.GlobalClient.DB.SQLSelect(&res, sql)

	if err != nil {
		log.Println("ERROR: User: ", err)
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
		"message": "Logout Success.",
	})
}
