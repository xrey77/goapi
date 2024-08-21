package users

import (
	"fmt"

	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

func ValidateFullname(fname string) int {
	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT FullName FROM user WHERE FullName=?`
	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, fname)
	if err2 != nil {
		return 0
	}
	defer rows.Close()
	// // iterate over the rows
	var userx []models.TempUsers
	for rows.Next() {
		var xuser models.TempUsers
		// unmarshal the row object to user
		err3 := rows.Scan(&xuser.FullName)
		if err3 != nil {
			fmt.Println(err3)
		}
		userx = append(userx, xuser)
	}
	rows.Close()
	if len(userx) == 0 {
		return 0
	}
	return 1
}

func ValidateEmail(mail string) int {
	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT Email FROM user WHERE Email=?`
	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, mail)
	if err2 != nil {
		return 0
	}
	defer rows.Close()
	// // iterate over the rows
	var userx []models.TempUsers
	for rows.Next() {
		var xuser models.TempUsers
		// unmarshal the row object to user
		err3 := rows.Scan(&xuser.Email)
		if err3 != nil {
			fmt.Println(err3)
		}
		userx = append(userx, xuser)
	}
	rows.Close()
	if len(userx) == 0 {
		return 0
	}
	return 1
}

func ValidateUsername(usrname string) int {
	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT UserName FROM user WHERE UserName=?`
	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, usrname)
	if err2 != nil {
		return 0
	}
	defer rows.Close()
	// // iterate over the rows
	var userx []models.TempUsers
	for rows.Next() {
		var xuser models.TempUsers
		// unmarshal the row object to user
		err3 := rows.Scan(&xuser.UserName)
		if err3 != nil {
			fmt.Println(err3)
		}
		userx = append(userx, xuser)
	}
	rows.Close()
	if len(userx) == 0 {
		return 0
	}
	return 1
}

func ValidateMailToken(mtoken string) int {
	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT MailToken FROM user WHERE MailToken=?`
	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, mtoken)
	if err2 != nil {
		return 0
	}
	defer rows.Close()
	// // iterate over the rows
	var userx []models.TempUsers
	for rows.Next() {
		var xuser models.TempUsers
		// unmarshal the row object to user
		err3 := rows.Scan(&xuser.MailToken)
		if err3 != nil {
			fmt.Println(err3)
		}
		userx = append(userx, xuser)
	}
	rows.Close()
	if len(userx) == 0 {
		return 0
	}
	return 1
}
