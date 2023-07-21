package controllers

import (
	"fmt"
	"time"

	"github.com/buger/jsonparser"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nikitamirzani323/isb_landingpage_api/models"
)

const Fielddomain_home_redis = "LISTDOMAIN_BACKEND"

func Domainhome(c *fiber.Ctx) error {
	var obj entities.Model_domain
	var arraobj []entities.Model_domain
	render_page := time.Now()
	resultredis, flag := helpers.GetRedis(Fielddomain_home_redis)
	jsonredis := []byte(resultredis)
	record_RD, _, _, _ := jsonparser.Get(jsonredis, "record")
	jsonparser.ArrayEach(record_RD, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		domain_id, _ := jsonparser.GetInt(value, "domain_id")
		domain_name, _ := jsonparser.GetString(value, "domain_name")
		domain_status, _ := jsonparser.GetString(value, "domain_status")
		domain_create, _ := jsonparser.GetString(value, "domain_create")
		domain_update, _ := jsonparser.GetString(value, "domain_update")

		obj.Domain_id = int(domain_id)
		obj.Domain_name = domain_name
		obj.Domain_status = domain_status
		obj.Domain_create = domain_create
		obj.Domain_update = domain_update
		arraobj = append(arraobj, obj)
	})

	if !flag {
		result, err := models.Fetch_domainHome()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": err.Error(),
				"record":  nil,
			})
		}
		helpers.SetRedis(Fielddomain_home_redis, result, 60*time.Minute)
		fmt.Println("DOMAIN MYSQL")
		return c.JSON(result)
	} else {
		fmt.Println("DOMAIN CACHE")
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Success",
			"record":  arraobj,
			"time":    time.Since(render_page).String(),
		})
	}
}
func DomainSave(c *fiber.Ctx) error {
	var errors []*helpers.ErrorResponse
	client := new(entities.Controller_domainsave)
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

	result, err := models.Save_domain(
		client_admin,
		client.Domain_name, client.Domain_status, client.Sdata, client.Domain_id)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": err.Error(),
			"record":  nil,
		})
	}

	_deleteredis_domain()
	return c.JSON(result)
}
func _deleteredis_domain() {
	val_master := helpers.DeleteRedis(Fielddomain_home_redis)
	fmt.Printf("Redis Delete BACKEND DOMAIN : %d", val_master)

}
