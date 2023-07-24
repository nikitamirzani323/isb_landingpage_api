package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nikitamirzani323/isb_landingpage_api/models"
)

const Fielddomain_home_redis = "LISTDOMAIN_FRONTEND_LANDINGPAGE"

func Domainhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_domain)
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

	client_origin := c.Request().Body()
	data_origin := []byte(client_origin)
	hostname, _ := jsonparser.GetString(data_origin, "client_hostname")
	fmt.Println("Request Body : ", string(data_origin))
	fmt.Println("BANNER CLIENT origin : ", hostname)

	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fielddomain_home_redis)
	jsonredis := []byte(resultredis)
	message_RD, _ := jsonparser.GetString(jsonredis, "message")
	domain_RD, _ := jsonparser.GetString(jsonredis, "domain")

	if !flag {
		result, err := models.Get_AllDomain()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fielddomain_home_redis, result, 60*time.Hour)
		fmt.Println("DOMAIN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("DOMAIN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": message_RD,
			"domain":  domain_RD,
			"time":    time.Since(render_page).String(),
		})
	}
}
