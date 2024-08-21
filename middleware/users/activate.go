package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err1)
	}

	db := middleware.CreateConnection()
	defer db.Close()

	res, errs2 := db.Exec("UPDATE user SET isactivated=? WHERE id=?", 1, id)
	if errs2 != nil {
		fmt.Println(errs2)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {

		// create the select sql query
		sqlStatement := `SELECT Email FROM user WHERE Id=?`
		// execute the sql statement
		rows, err2 := db.Query(sqlStatement, id)
		if err2 != nil {
			msg := map[string]string{"message": err2.Error()}
			json.NewEncoder(w).Encode(msg)
			return
		}
		defer rows.Close()
		id = 0
		// GET USER EMAIL
		var usermail []models.TempUsers
		for rows.Next() {
			var xmail models.TempUsers

			// unmarshal the row object to user
			err3 := rows.Scan(&xmail.Email)
			if err3 != nil {
				msg := map[string]string{"message": err3.Error()}
				json.NewEncoder(w).Encode(msg)
			}
			usermail = append(usermail, xmail)
		}
		var myEmail = usermail[0].Email

		//SEND EMAIL FOR ACTIVATION
		var msgbody string = "<div>Congrantiolations !, your account has been activated.</div><br/></br><div>Best Regards,<br/>Administrator</div>"
		activate := ActivateAccount(msgbody, "Account activation confirmation", myEmail)
		if activate == 0 {
			fmt.Println("not sent..")
		}
		//DISPLAY MESSAGE TO THE BROWSER
		AddMessage := "<div style=\"text-align: center;margin-top: 200px;font-size: 38px;font-weight: bolder;\">Congratulations! your account has been activated.</div>"
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, AddMessage)
		return

	}
}
