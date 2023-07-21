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

const database_uom_local = configs.DB_tbl_mst_uom

func Fetch_uomHome(search string, page int) (helpers.Responsepaging, error) {
	var obj entities.Model_uom
	var arraobj []entities.Model_uom
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
	sql_selectcount += "COUNT(iduom) as totaluom  "
	sql_selectcount += "FROM " + database_uom_local + "  "
	if search != "" {
		sql_selectcount += "WHERE LOWER(nmuom) LIKE '%" + strings.ToLower(search) + "%' "
		sql_selectcount += "OR LOWER(nmuom) LIKE '%" + strings.ToLower(search) + "%' "
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
	sql_select += "iduom , nmuom,  statusuom, "
	sql_select += "createuom, to_char(COALESCE(createdateuom,now()), 'YYYY-MM-DD HH24:ii:ss') as createdateuom, "
	sql_select += "updateuom, to_char(COALESCE(updatedateuom,now()), 'YYYY-MM-DD HH24:ii:ss') as updatedateuom "
	sql_select += "FROM " + database_uom_local + " "
	if search == "" {
		sql_select += "ORDER BY updatedateuom DESC  OFFSET " + strconv.Itoa(offset) + " LIMIT " + strconv.Itoa(perpage)
	} else {
		sql_select += "WHERE LOWER(nmuom) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "OR LOWER(nmuom) LIKE '%" + strings.ToLower(search) + "%' "
		sql_select += "ORDER BY updatedateuom DESC  LIMIT " + strconv.Itoa(perpage)
	}

	row, err := con.QueryContext(ctx, sql_select)
	helpers.ErrorCheck(err)
	for row.Next() {
		var (
			iduom_db                                                       int
			nmuom_db, statusuom_db                                         string
			createuom_db, createdateuom_db, updateuom_db, updatedateuom_db string
		)

		err = row.Scan(&iduom_db, &nmuom_db, &statusuom_db,
			&createuom_db, &createdateuom_db, &updateuom_db, &updatedateuom_db)

		helpers.ErrorCheck(err)
		create := ""
		update := ""
		if createuom_db != "" {
			create = createuom_db + ", " + createdateuom_db
		}
		if updateuom_db != "" {
			update = updateuom_db + ", " + updatedateuom_db
		}

		obj.Uom_id = iduom_db
		obj.Uom_name = nmuom_db
		obj.Uom_status = statusuom_db
		obj.Uom_create = create
		obj.Uom_update = update
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
func Save_uom(admin, name, status, sData string, idrecord int) (helpers.Response, error) {
	var res helpers.Response
	msg := "Failed"
	tglnow, _ := goment.New()
	render_page := time.Now()
	if sData == "New" {
		sql_insert := `
				insert into
				` + database_uom_local + ` (
					iduom , nmuom, statusuom, 
					createuom, createdateuom
				) values (
					$1, $2, $3,  
					$4, $5 
				)
			`
		field_column := database_uom_local + tglnow.Format("YYYY")
		idrecord_counter := Get_counter(field_column)
		flag_insert, msg_insert := Exec_SQL(sql_insert, database_uom_local, "INSERT",
			tglnow.Format("YY")+strconv.Itoa(idrecord_counter), name, status,
			admin, tglnow.Format("YYYY-MM-DD HH:mm:ss"))

		if flag_insert {
			msg = "Succes"
		} else {
			fmt.Println(msg_insert)
		}
	} else {
		sql_update := `
				UPDATE 
				` + database_uom_local + `  
				SET nmuom =$1, statusuom =$2, 
				updatecurr=$3, updatedatecurr=$4   
				WHERE idcurr=$5 
			`

		flag_update, msg_update := Exec_SQL(sql_update, database_uom_local, "UPDATE",
			name, status, admin,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"), idrecord)

		if flag_update {
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
