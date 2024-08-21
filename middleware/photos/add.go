// http://localhost:5100/photo/add

// {
// 	"phototitle": "Ang Dating Daan",
// 	"photoimage": "datingdaan.png"
//   }

package photos

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

//ADD NEW PHOTO
//DISPLAY BASE54 BINAY IN FRONT END
//data:image/jpeg;base64,ASDFASDF)
func Add(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json") // for json data
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded") //for form data

	//AUTHORIZATION==========
	_, errAuth := middleware.ExtractTokenID(r)
	if errAuth != nil {
		msg := map[string]string{"message": "UnAuthorized Access."}
		json.NewEncoder(w).Encode(msg)
		return
	}
	//========================

	//CONNECT TO MYSQL DATABASE
	db := middleware.CreateConnection()
	defer db.Close()

	//GET FORM INPUT DATA
	title := r.FormValue("phototitle")

	//START GET NEXT RECORD NO=======
	sql := `SELECT MAX(id)+1 FROM photos ORDER BY id`
	// execute the sql statement
	xrows, _ := db.Query(sql)
	defer xrows.Close()

	// GET USER EMAIL
	var photo1 []models.TempPhotos
	for xrows.Next() {
		var photo2 models.TempPhotos
		// unmarshal the row object to user
		err3 := xrows.Scan(&photo2.ID)
		if err3 != nil {
			msg := map[string]string{"message": err3.Error()}
			json.NewEncoder(w).Encode(msg)
		}
		photo1 = append(photo1, photo2)
	}
	var myIDno = strconv.FormatInt(photo1[0].ID, 10)
	//END============================

	//START UPLOAD IMAGE=======
	imageName, err := FileUpload(r, myIDno)
	if err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	//GET IMAGE EXTENSION AND ASSIGN NEW FILENAME
	var pic string
	img := strings.ToLower(imageName)
	if strings.Contains(img, ".jpg") {
		pic = myIDno + ".jpg"
	} else if strings.Contains(img, ".jpeg") {
		pic = myIDno + ".jpeg"
	} else if strings.Contains(img, ".png") {
		pic = myIDno + ".png"
	} else if strings.Contains(img, ".gif") {
		pic = myIDno + ".gif"
	}

	//END======================

	//START CONVERT TO BASE54 BINARY TO STRING======
	// var pic string = base64.StdEncoding.EncodeToString([]byte(data))
	//END===========================================

	var val = ValidateTitle(title)
	if val == 1 {
		msg := map[string]string{"message": "Photo Title already taken."}
		json.NewEncoder(w).Encode(msg)
		return
	}
	//INSERT
	res, errs2 := db.Exec("INSERT INTO photos(photo_title,photo_image) VALUES(?,?)", &title, &pic)
	if errs2 != nil {
		fmt.Println("err2", errs2)
		return
	}
	pid, err3 := res.LastInsertId()
	if err3 != nil {
		log.Fatal("err3", err3)
		return
	}
	id := strconv.FormatInt(pid, 10)
	msg := map[string]string{"message": "New Photo ID No. " + id + " has been created."}
	json.NewEncoder(w).Encode(msg)
}

func FileUpload(r *http.Request, fstr string) (string, error) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	r.ParseMultipartForm(32 << 20)
	//ParseMultipartForm parses a request body as multipart/form-data
	file, handler, err := r.FormFile("photoimage") //retrieve the file from form data
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
		newFile = fstr + ".jpg"
	} else if strings.Contains(img, ".jpeg") {
		newFile = fstr + ".jpeg"
	} else if strings.Contains(img, ".png") {
		newFile = fstr + ".png"
	} else if strings.Contains(img, ".gif") {
		newFile = fstr + ".gif"
	}

	f, err := os.OpenFile("public/images/photos/"+newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	io.Copy(f, file)
	//here we save our file to our path
	return handler.Filename, nil
}
