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

const Fielduom_home_redis = "LISTUOM_BACKEND"
const Fielduom_frontend_redis = "LISTUOM_FRONTEND"

func Uomhome(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_uom)
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
	if client.Uom_search != "" {
		val_uom := helpers.DeleteRedis(Fielduom_home_redis + "_" + strconv.Itoa(client.Uom_page) + "_" + client.Uom_search)
		fmt.Printf("Redis Delete BACKEND UOM : %d", val_uom)
	}

	var obj entities.Model_uom
	var arraobj []entities.Model_uom
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fielduom_home_redis + "_" + strconv.Itoa(client.Uom_page) + "_" + client.Uom_search)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		uom_id, _ := jsonparser.GetInt(value, "uom_id")
		uom_name, _ := jsonparser.GetString(value, "uom_name")
		uom_status, _ := jsonparser.GetString(value, "uom_status")
		uom_create, _ := jsonparser.GetString(value, "uom_create")
		uom_update, _ := jsonparser.GetString(value, "uom_update")

		obj.Uom_id = int(uom_id)
		obj.Uom_name = uom_name
		obj.Uom_status = uom_status
		obj.Uom_create = uom_create
		obj.Uom_update = uom_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_uomHome(client.Uom_search, client.Uom_page)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fielduom_home_redis+"_"+strconv.Itoa(client.Uom_page)+"_"+client.Uom_search, result, 1*time.Hour)
		fmt.Println("UOM MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("UOM CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func UomSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_uomsave)
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

	result, err := models.Save_uom(
		client_admin,
		client.Uom_name, client.Uom_status, client.Sdata, client.Uom_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_uom(client.Uom_search, client.Uom_page)
	return c.JSON(result)
}
func _deleteredis_uom(search string, page int) {
	val_master := helpers.DeleteRedis(Fielduom_home_redis + "_" + strconv.Itoa(page) + "_" + search)
	log.Printf("Redis Delete BACKEND UOM : %d", val_master)

	val_client := helpers.DeleteRedis(Fielduom_frontend_redis)
	log.Printf("Redis Delete FRONTEND UOM : %d", val_client)

}
