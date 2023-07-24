package models

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/isb_landingpage_api/config"
	"github.com/nikitamirzani323/isb_landingpage_api/db"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
)

func Get_AllDomain() (helpers.ResponseDomain, error) {
	var res helpers.ResponseDomain
	msg := "Data Not Found"
	ctx := context.Background()
	con := db.CreateCon()
	domain_client := ""

	sql_select := `SELECT
		nmdomain   
		FROM ` + config.DB_tbl_mst_domain + `  
		WHERE statusdomain = 'RUNNING' 
		LIMIT 1 
	`
	row, err := con.QueryContext(ctx, sql_select)
	defer row.Close()
	helpers.ErrorCheck(err)

	for row.Next() {
		var nmdomain_db string
		err = row.Scan(&nmdomain_db)
		if err != nil {
			return res, err
		}
		domain_client = nmdomain_db
		msg = "Success"
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Domain = domain_client

	return res, nil
}
