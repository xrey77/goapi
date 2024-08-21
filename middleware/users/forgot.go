package users

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//GET PARAMETERS
	params := mux.Vars(r)
	email := params["email"]

	//IF EMAIL NOT EXISTS
	if ValidateEmail(email) == 0 {
		msg := map[string]string{"message": "Email Address not found."}
		json.NewEncoder(w).Encode(msg)
	}

	//CONNECT TO MYSQL DATABASE
	db := middleware.CreateConnection()
	defer db.Close()

	// EXECUTE SQL STATEMENT
	sqlStatement := `SELECT MailToken FROM user WHERE Email=?`
	erows, err2 := db.Query(sqlStatement, email)
	if err2 != nil {
		msg := map[string]string{"message": err2.Error()}
		json.NewEncoder(w).Encode(msg)
		return
	}
	defer erows.Close()

	// GET EXISTING USER E-MAIL TOKEN
	var usermail []models.TempUsers
	for erows.Next() {
		var xmail models.TempUsers

		// unmarshal the row object to user
		err3 := erows.Scan(&xmail.MailToken)
		if err3 != nil {
			msg := map[string]string{"message": err3.Error()}
			json.NewEncoder(w).Encode(msg)
		}
		usermail = append(usermail, xmail)
	}
	myEmailToken := usermail[0].MailToken

	//GENERATE E-MAIL TOKEN
	vtoken := strconv.Itoa(xtokens(int(myEmailToken), 999999))

	//UPDATE NEW E-MAIL TOKEN FROM DATABASE
	_, errs2 := db.Exec("UPDATE user SET MailToken=? WHERE Email=?", vtoken, email)
	if errs2 != nil {
		fmt.Println(errs2)
	}

	//SEND MAIL TOKEN
	var msgbody string = "<div>Copy and Paste E-Mail Token <span style=\"font-size:30px;\" >" + vtoken + "</span> , for you to change your password.</div>"
	activate := ActivateAccount(msgbody, "E-Mail Token", email)
	if activate > 0 {
		msg := map[string]string{"message": "E-Mail Token has been sent to " + email}
		json.NewEncoder(w).Encode(msg)
	}
}

func SendToken(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	totp := params["token"]

	if ValidateMailToken(totp) == 0 {
		msg := map[string]string{"message": "E-Mail Token does not match. "}
		json.NewEncoder(w).Encode(msg)
		return
	}
	msg := map[string]string{"message": "E-Mail Token has been sent. "}
	json.NewEncoder(w).Encode(msg)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginmoodel models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginmoodel)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	if ValidateUsername(loginmoodel.UserName) == 0 {
		msg := map[string]string{"message": "Username not found."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	//HASH NEW PASSWORD
	pwd := hashAndSalt([]byte(loginmoodel.PassWord))
	fmt.Println(pwd)

	//UPDATE PASSWORD
	db := middleware.CreateConnection()
	defer db.Close()

	//UPDATE
	res, errs2 := db.Exec("UPDATE user SET Password=? WHERE UserName=?", pwd, loginmoodel.UserName)
	if errs2 != nil {
		fmt.Println(errs2)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		msg := map[string]string{"message": "Password has been updated successfully."}
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := map[string]string{"message": "Unable to update password."}
		json.NewEncoder(w).Encode(msg)
	}

}

func xtokens(low, hi int) int {
	return low + rand.Intn(hi-low)
}
