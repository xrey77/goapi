// http://localhost:5100/photo/update/100
// {
// 	"phototitle": "Bread of LIfe Ministries",
// 	"photoimage": "lionheart.png"
//   }

package photos

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

//UPDATE PHOTO
func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//AUTHORIZATION
	_, errAuth := middleware.ExtractTokenID(r)
	if errAuth != nil {
		msg := map[string]string{"message": "UnAuthorized Access."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	params := mux.Vars(r)
	id, errId := strconv.Atoi(params["id"])
	if errId != nil {
		log.Fatalf("Unable to convert the string into int.  %v", errId)
	}

	var modelphoto models.Photo
	// decode the json request to contacts
	err := json.NewDecoder(r.Body).Decode(&modelphoto)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	db := middleware.CreateConnection()
	defer db.Close()

	//UPDATE
	res, errs2 := db.Exec("UPDATE photos SET photo_title=?, photo_image=? WHERE id=?", modelphoto.Photo_title, modelphoto.Photo_image, id)
	if errs2 != nil {
		fmt.Println(errs2)
	}
	xid := params["id"]
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		msg := map[string]string{"message": "Photo ID No." + xid + " Updated Successfully."}
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := map[string]string{"message": "Photo ID No. " + xid + " does not exists."}
		json.NewEncoder(w).Encode(msg)
	}

}
