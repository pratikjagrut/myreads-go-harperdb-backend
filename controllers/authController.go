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

	res, err := database.GlobalClient.Where("users", "email", user.Email)
	if err != nil {
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: Register: %v", err),
		})
	}

	if len(res) != 0 {
		c.Status(fiber.StatusFound)
		return c.JSON(fiber.Map{
			"message": "This email id is already registered.",
		})
	}

	var m = map[string]string{
		"email":    user.Email,
		"name":     user.Name,
		"password": string(user.Password),
	}

	_, err = database.GlobalClient.Insert("users", m)

	if err != nil {
		c.Status(fiber.StatusExpectationFailed)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("ERROR: Register: %v", err),
		})
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	res, err := database.GlobalClient.Where("users", "email", data["email"])
	if err != nil {
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

	res, err := database.GlobalClient.Where("users", "id", claims.Issuer)

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
