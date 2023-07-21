package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nikitamirzani323/isb_landingpage_api/configs"
	"github.com/nikitamirzani323/isb_landingpage_api/db"
	"github.com/nikitamirzani323/isb_landingpage_api/entities"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nleeper/goment"
)

const database_curr_local = configs.DB_tbl_mst_currency

func Fetch_currencyHome(search string, page int) (helpers.Responsepaging, error) {
	var obj entities.Model_currency
	var arraobj []entities.Model_currency
	var res helpers.Responsepaging
	msg := "Data Not Found"
	con := db.CreateCon()
	ctx := context.Background()
	start := time.Now()

	perpage := 50
	totalrecord := 0
	offset := page
	sql_selectcount := ""
	sql_selectcount += ""
	sql_selectcount += "SELECT "
	sql_selectcount += "COUNT(idcurr) as totalcurr  "
	sql_selectcount += "FROM " + database_curr_local + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(nmcurr) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nmcurr) LIKE '%" + strings.ToLower(search) + "%' "
	}

	row_selectcount := con.QueryRowContext(ctx, sql_selectcount)
	switch e_selectcount := row_selectcount.Scan(&totalrecord); e_selectcount {
	case sql.ErrNoRows:
	case nil:
	default:
		helpers.ErrorCheck(e_selectcount)
	}

	sql_select := ""
	sql_select += ""
	sql_select += "SELECT "
	sql_select += "idcurr , nmcurr, "
	sql_select += "createcurr, to_char(COALESCE(createdatecurr,now()), 'YYYY-MM-DD HH24:ii:ss') as createdatecurr, "
	sql_select += "updatecurr, to_char(COALESCE(updatedatecurr,now()), 'YYYY-MM-DD HH24:ii:ss') as updatedatecurr "
	sql_select += "FROM " + database_uom_local + " "
	if search == "" {
		sql_select += "ORDER BY updatecurr DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(nmcurr) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(nmcurr) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY updatecurr DESC  LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			idcurr_db, nmcurr_db                                               string
			createcurr_db, createdatecurr_db, updatecurr_db, updatedatecurr_db string
		)

		err = row.Scan(&idcurr_db, &nmcurr_db,
			&createcurr_db, &createdatecurr_db, &updatecurr_db, &updatedatecurr_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createcurr_db != "" {
			create = createcurr_db + ", " + createdatecurr_db
		}
		if updatecurr_db != "" {
			update = updatecurr_db + ", " + updatedatecurr_db
		}

		obj.Currency_id = idcurr_db
		obj.Currency_name = nmcurr_db
		obj.Currency_create = create
		obj.Currency_update = update
		arraobj = append(arraobj, obj)
		msg = "Success"
	}
	defer row.Close()

	res.Status = fiber.StatusOK
	res.Message = msg
	res.Record = arraobj
	res.Perpage = perpage
	res.Totalrecord = totalrecord
	res.Time = time.Since(start).String()

	return res, nil
}
func Save_currency(admin, name, idrecord, sData string) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	flag := false

	if sData == "New" {
		flag = CheckDB(database_curr_local, "idcurr", idrecord)
		if !flag {
			sql_insert := `
				insert into
				` + database_curr_local + ` (
					idcurr , nmcurr, 
					createcurr, createdatecurr
				) values (
					$1, $2, 
					$3, $4 
				)
			`

			flag_insert, msg_insert := Exec_SQL(sql_insert, database_curr_local, "INSERT",
				idrecord, name,
				admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

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
				` + database_curr_local + `  
				SET nmcurr =$1, 
				updatecurr=$2, updatedatecurr=$3  
				WHERE idcurr=$4 
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_curr_local, "UPDATE",
			name, admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

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
