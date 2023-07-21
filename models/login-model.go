package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/nikitamirzani323/isb_landingpage_api/configs"
	"github.com/nikitamirzani323/isb_landingpage_api/db"
	"github.com/nikitamirzani323/isb_landingpage_api/helpers"
	"github.com/nleeper/goment"
)

func Login_Model(username, password, ipaddress, timezone string) (bool, string, error) {
	con := db.CreateCon()
	ctx := context.Background()
	flag := false
	tglnow, _ := goment.New()
	var passwordDB, idadminDB string
	sql_select := `
			SELECT
			password, idadmin    
			FROM ` + configs.DB_tbl_admin + ` 
			WHERE username  = $1
			AND statuslogin = 'Y' 
		`

	fmt.Println(sql_select, username)
	row := con.QueryRowContext(ctx, sql_select, username)
	switch e := row.Scan(&passwordDB, &idadminDB); e {
	case sql.ErrNoRows:
		return false, "", errors.New("Username and Password Not Found")
	case nil:
		flag = true
	default:
		return false, "", errors.New("Username and Password Not Found")
	}

	hashpass := helpers.HashPasswordMD5(password)

	if hashpass != passwordDB {
		return false, "", nil
	}

	if flag {
		sql_update := `
			UPDATE ` + configs.DB_tbl_admin + ` 
			SET lastlogin=$1, ipaddress=$2 , timezone=$3, 
			updateadmin=$4,  updatedateadmin=$5  
			WHERE username  = $6 
			AND statuslogin = 'Y' 
		`
		flag_update, msg_update := Exec_SQL(sql_update, configs.DB_tbl_admin, "UPDATE",
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			ipaddress,
			timezone,
			username,
			tglnow.Format("YYYY-MM-DD HH:mm:ss"),
			username)

		if flag_update {
			flag = true
			log.Println(msg_update)
		} else {
			log.Println(msg_update)
		}
	}

	return true, idadminDB, nil
}
