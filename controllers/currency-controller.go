package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nikitamirzani323/isb_landingpage_api/models"
)

const Fieldcurrency_home_redis = "LISTCURRECNY_BACKEND"
const Fieldcurrency_frontend_redis = "LISTCURRENCY_FRONTEND"

func Currencyhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_currency)
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
	if client.Currency_search != "" {
		val_curr := helpers.DeleteRedis(Fieldcurrency_home_redis + "_" + strconv.Itoa(client.Currency_page) + "_" + client.Currency_search)
		fmt.Printf("Redis Delete BACKEND CURR : %d", val_curr)
	}

	var obj entities.Model_currency
	var arraobj []entities.Model_currency
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fielduom_home_redis + "_" + strconv.Itoa(client.Currency_page) + "_" + client.Currency_search)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		currency_id, _ := jsonparser.GetString(value, "currency_id")
		currency_name, _ := jsonparser.GetString(value, "currency_name")
		currency_create, _ := jsonparser.GetString(value, "currency_create")
		currency_update, _ := jsonparser.GetString(value, "currency_update")

		obj.Currency_id = currency_id
		obj.Currency_name = currency_name
		obj.Currency_create = currency_create
		obj.Currency_update = currency_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_currencyHome(client.Currency_search, client.Currency_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fieldcurrency_home_redis+"_"+strconv.Itoa(client.Currency_page)+"_"+client.Currency_search, result, 1*time.Hour)
		fmt.Println("CURRENCY MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("CURRENCY CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func CurrencySave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_currencysave)
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
	user := c.Locals("jwt").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	temp_decp := helpers.Decryption(name)
	client_admin, _ := helpers.Parsing_Decry(temp_decp, "==")

	result, err := models.Save_currency(
		client_admin,
		client.Currency_name, client.Currency_id, client.Sdata)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_currency(client.Currency_search, client.Currency_page)
	return c.JSON(result)
}
func _deleteredis_currency(search string, page int) {
	val_master := helpers.DeleteRedis(Fieldcurrency_home_redis + "_" + strconv.Itoa(page) + "_" + search)
	log.Printf("Redis Delete BACKEND CURRENCY : %d", val_master)

	val_client := helpers.DeleteRedis(Fieldcurrency_frontend_redis)
	log.Printf("Redis Delete FRONTEND CURRENCY : %d", val_client)

}
