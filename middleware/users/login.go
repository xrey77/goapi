// http://localhost:5100/users/login
// {
// 	"username": "Rey",
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
	u "github.com/reynald/goapi/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.UserLogin

	// decode the json request to user
	err1 := json.NewDecoder(r.Body).Decode(&user)
	if err1 != nil {
		log.Fatalf("Unable to decode the request body.  %v", err1)
	}

	uname := ValidateUsername(user.UserName)
	if uname == 0 {
		msg := map[string]string{"message": "Username not found."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := "SELECT id, UserName, Password, FullName,MobileNO,Email,isactivated FROM user WHERE UserName = ?"
	row, err2 := db.Query(sqlStatement, user.UserName)
	if err2 != nil {
		fmt.Println(err2)
	}
	defer row.Close()
	var juser models.JsonUser
	for row.Next() {
		var xuser models.TempUsers

		err3 := row.Scan(&xuser.ID, &xuser.UserName, &xuser.Password, &xuser.FullName, &xuser.MobileNo, &xuser.Email, &xuser.Isactivated)
		if err3 != nil {
			msg := map[string]string{"message": "Credentials not found."}
			json.NewEncoder(w).Encode(msg)
			log.Println("error", err3)
			return
		} else {
			juser.ID = xuser.ID
			juser.FullName = xuser.FullName
			juser.Email = xuser.Email
			juser.MobileNo = xuser.MobileNo
			juser.UserName = xuser.UserName

			//VALIDATE IF USER IS ACTIVATED
			if xuser.Isactivated == 0 {
				msg := map[string]string{"message": "Please activate your account."}
				json.NewEncoder(w).Encode(msg)
				return
			}

			//VALIDATE IF USER HAS BEEN BLOCKED
			if xuser.Isactivated > 1 {
				msg := map[string]string{"message": "Your account has been deactivated, please contact the Administrator."}
				json.NewEncoder(w).Encode(msg)
				return
			}

			if comparePassword(xuser.Password, getPassword(user.PassWord)) {

				//Create JWT token
				// tk := &models.Token{UserId: user.ID}
				// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
				// tokenString, _ := token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

				//DELETE USERID COOKIE
				// cokie := http.Cookie{Name: "USERID", MaxAge: -1}
				// http.SetCookie(w, &cokie)

				//CREATE USERID COOKIE
				UID := strconv.FormatInt(juser.ID, 10)
				c1 := http.Cookie{Name: "USERID", Path: "/", Value: UID, MaxAge: 3600}
				http.SetCookie(w, &c1)
				UID = ""

				tokenString, _ := middleware.CreateToken(uint32(juser.ID))
				user.Token = tokenString //Store the token in the response
				juser.Token = user.Token

				w.Header().Set("Content-Type", "application/json")
				prettyJSON, err := json.MarshalIndent(juser, "", "    ")
				if err != nil {
					msg := map[string]string{"Failed to generate json": err.Error()}
					json.NewEncoder(w).Encode(msg)
				}
				resp := u.Message(true, "Logged In")
				resp["user"] = user

				w.Write(prettyJSON)

			} else {
				msg := map[string]string{"message": "Access Denied."}
				json.NewEncoder(w).Encode(msg)
			}

		}
	}
}

func comparePassword(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func getPassword(pwd string) []byte {
	return []byte(pwd)
}
