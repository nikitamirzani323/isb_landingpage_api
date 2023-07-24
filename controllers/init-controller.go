package controllers

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
)

func CheckLogin(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_clientinit)
	validate := validator.New()
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}
	err := validate.Struct(client)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element helpers.ErrorResponse
			element.Field = err.StructField()
			element.Tag = err.Tag()
			errors = append(errors, &element)
		}
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "validation",
			"record":  errors,
		})
	}

	dataclient := client.Client_hostname
	dataclient_encr, keymap := helpers.Encryption(dataclient)
	dataclient_encr_final := dataclient_encr + "|" + strconv.Itoa(keymap)
	t, err := helpers.GenerateNewAccessToken(dataclient_encr_final)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"token":  t,
	})
}
func Home(c *fiber.Ctx) error {
	client := new(entities.Controller_clientinit)
	if err := c.BodyParser(client); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	fmt.Println("name", name)
	fmt.Println("Decp", temp_decp)

	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success",
		"record":  nil,
	})
}
