package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/reynald/goapi/middleware"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded") //for form data

	_id := r.FormValue("idno")
	_fullname := r.FormValue("fullname")
	_email := r.FormValue("email")
	_mobileno := r.FormValue("mobileno")
	_password := r.FormValue("password")

	//START UPLOAD IMAGE=======
	imageName, err := ProfileUpload(r, _id)
	if err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	//GET IMAGE EXTENSION AND ASSIGN NEW FILENAME
	var _pic string
	img := strings.ToLower(imageName)
	if strings.Contains(img, ".jpg") {
		_pic = "000" + _id + ".jpg"
	} else if strings.Contains(img, ".jpeg") {
		_pic = "000" + _id + ".jpeg"
	} else if strings.Contains(img, ".png") {
		_pic = "000" + _id + ".png"
	} else if strings.Contains(img, ".gif") {
		_pic = "000" + _id + ".gif"
	}

	//HASH NEW PASSWORD
	_pwd := hashAndSalt([]byte(_password))

	//UPDATE PASSWORD
	db := middleware.CreateConnection()
	defer db.Close()

	//UPDATE
	res, errs2 := db.Exec("UPDATE user SET FullName=?, Email=?, MobileNo=?, Password=?, UserPic=? WHERE id=?", _fullname, _email, _mobileno, _pwd, _pic, _id)
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

func ProfileUpload(r *http.Request, fstr string) (string, error) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	r.ParseMultipartForm(32 << 20)
	//ParseMultipartForm parses a request body as multipart/form-data
	file, handler, err := r.FormFile("userpic") //retrieve the file from form data
	//replace file with the key your sent your image with
	if err != nil {
		return "", err
	}

	defer file.Close() //close the file when we finish
	//this is path which  we want to store the file
	// f, err := os.OpenFile("public/images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	var newFile string
	img := strings.ToLower(handler.Filename)
	if strings.Contains(img, ".jpg") {
		newFile = "000" + fstr + ".jpg"
	} else if strings.Contains(img, ".jpeg") {
		newFile = "000" + fstr + ".jpeg"
	} else if strings.Contains(img, ".png") {
		newFile = "000" + fstr + ".png"
	} else if strings.Contains(img, ".gif") {
		newFile = "000" + fstr + ".gif"
	}

	f, err := os.OpenFile("public/images/users/"+newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	io.Copy(f, file)
	//here we save our file to our path
	return handler.Filename, nil
}
