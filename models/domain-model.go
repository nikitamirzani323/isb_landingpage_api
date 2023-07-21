package models

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/isb_landingpage_api/configs"
	"github.com/nikitamirzani323/isb_landingpage_api/db"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nleeper/goment"
)

const database_domain_local = configs.DB_tbl_mst_domain

func Fetch_domainHome() (helpers.Response, error) {
	var obj entities.Model_domain
	var arraobj []entities.Model_domain
	var res helpers.Response
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	sql_select := `SELECT 
			iddomain , nmdomain, statusdomain,  
			createdomain, to_char(COALESCE(createdatedomain,now()), 'YYYY-MM-DD HH24:MI:SS'), 
			updatedomain, to_char(COALESCE(updatedatedomain,now()), 'YYYY-MM-DD HH24:MI:SS') 
			FROM ` + database_domain_local + `  
			ORDER BY iddomain DESC   `

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iddomain_db                                                                int
			nmdomain_db, statusdomain_db                                               string
			createdomain_db, createdatedomain_db, updatedomain_db, updatedatedomain_db string
		)

		err = row.Scan(&iddomain_db, &nmdomain_db, &statusdomain_db,
			&createdomain_db, &createdatedomain_db, &updatedomain_db, &updatedatedomain_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createdomain_db != "" {
			create = createdomain_db + ", " + createdatedomain_db
		}
		if updatedomain_db != "" {
			update = updatedomain_db + ", " + updatedatedomain_db
		}

		obj.Domain_id = iddomain_db
		obj.Domain_name = nmdomain_db
		obj.Domain_status = statusdomain_db
		obj.Domain_create = create
		obj.Domain_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_domain(admin, nmdomain, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_domain_local, "nmdomain", nmdomain)
		if !flag {
			sql_insert := `
				insert into
				` + database_domain_local + ` (
					iddomain , nmdomain, statusdomain, 
					createdomain, createdatedomain
				) values (
					$1, $2, $3, 
					$4, $5
				)
			`
			field_column := database_domain_local + tglnow.Format("YYYY")
			idrecord_counter := Get_counter(field_column)
			flag_insert, msg_insert := Exec_SQL(sql_insert, database_domain_local, "INSERT",
				tglnow.Format("YY")+strconv.Itoa(idrecord_counter), nmdomain, status,
				admin,
				tglnow.Format("YYYY-MM-DD HH:mm:ss"))

			if flag_insert {
				msg = "Succes"
			} else {
				fmt.Println(msg_insert)
			}
		} else {
			msg = "Duplicate Entry"
		}
	} else {
		sql_update := `
				UPDATE 
				` + configs.DB_tbl_mst_domain + `  
				SET nmdomain =$1, statusdomain=$2, 
				updatedomain=$3, updatedatedomain=$4 
				WHERE iddomain =$5 
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_domain_local, "UPDATE",
			nmdomain, status,
			admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			idrecord)

		if flag_update {
			flag = true
			msg = "Succes"
		} else {
			fmt.Println(msg_update)
		}
	}

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = nil
	res.Time = time.Since(render_page).String()

	return res, nil
}
