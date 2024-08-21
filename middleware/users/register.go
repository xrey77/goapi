// http://localhost:5100/users/register
// {
// 	"full_name": "Rey Gragasin",
// 	"email": "reynald88@yahoo.com",
// 	"mobile_no": "343434",
// 	"username": "Reynald",
// 	"passwd": "Reynald88"
//   }

package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
	"golang.org/x/crypto/bcrypt"
)

//REGISTER NEW USER ACCOUNT
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	db := middleware.CreateConnection()
	defer db.Close()

	//CHECK AND BLOCK FULLNAME
	var valFullname = ValidateFullname(user.FullName)
	if valFullname == 1 {
		msg := map[string]string{"message": "Fullname already taken."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	//CHECK AND BLOCK EMAIL ADDRESS
	var valEmail = ValidateEmail(user.Email)
	if valEmail == 1 {
		msg := map[string]string{"message": "Email Address already taken."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	//CHECK AND BLOCK USERNAME
	var valUsername = ValidateUsername(user.UserName)
	if valUsername == 1 {
		msg := map[string]string{"message": "Username already taken."}
		json.NewEncoder(w).Encode(msg)
		return
	}
	xbyte := getPwd(user.Password)
	xhashPwd := hashAndSalt(xbyte)

	//INSERT
	sqlStatement := "INSERT INTO user(Fullname, Email, MobileNo, UserName, Password) VALUES(?, ?, ?, ?, ?)"
	res, errs := db.Exec(sqlStatement, &user.FullName, &user.Email, &user.MobileNo, &user.UserName, &xhashPwd)
	if errs != nil {
		msg := map[string]string{"message": errs.Error()}
		json.NewEncoder(w).Encode(msg)
		return
	}

	pid, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}

	//CREATE USER EMAIL COOKIE
	c2 := http.Cookie{Name: "USERMAIL", Path: "/", Value: user.Email, MaxAge: 3600}
	http.SetCookie(w, &c2)

	id := strconv.FormatInt(pid, 10)

	//SEND EMAIL FOR ACTIVATION
	var url string = "http://192.168.1.16:5100/api/users/activate/" + id
	var msgbody string = "<div>Please click button below to activate your account.</div><div><a style=\"background-color: green;color:white;font-size: 20px;border-radius: 20px;text-decoration: none;\" href=" + url + ">&nbsp;&nbsp;Activate Account&nbsp;&nbsp;</a></div>"
	var xmail string = user.Email
	activate := ActivateAccount(msgbody, "Activate your account.", xmail)
	if activate == 0 {
		fmt.Println("not sent..")
	}

	msg := map[string]string{"message": "Email has been sent to " + user.Email + ", please activate your account."}
	json.NewEncoder(w).Encode(msg)

}

func getPwd(pwd string) []byte {
	return []byte(pwd)
}

func hashAndSalt(pwd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
