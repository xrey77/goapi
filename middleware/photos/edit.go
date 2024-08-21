// http://localhost:5100/photo/edit/70

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

//EDIT/GET PHOTO BY ID
func Edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//AUTHORIZATION
	_, err := middleware.ExtractTokenID(r)
	if err != nil {
		msg := map[string]string{"message": "UnAuthorized Access."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	params := mux.Vars(r)
	id, err1 := strconv.Atoi(params["id"])
	if err1 != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err1)
	}

	db := middleware.CreateConnection()
	defer db.Close()
	var photox []models.TempPhotos

	// create the select sql query
	sqlStatement := `SELECT id, photo_title, photo_image FROM photos WHERE ID=? ORDER BY id`

	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, id)
	if err2 != nil {
		msg := map[string]string{"message": err2.Error()}
		json.NewEncoder(w).Encode(msg)
		return
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var xphoto models.TempPhotos

		// unmarshal the row object to user
		err3 := rows.Scan(&xphoto.ID, &xphoto.Photo_title, &xphoto.Photo_image)
		if err3 != nil {
			msg := map[string]string{"message": err3.Error()}
			json.NewEncoder(w).Encode(msg)
		}
		photox = append(photox, xphoto)
	}

	fmt.Println(photox)

	// sEnc := b64.StdEncoding.EncodeToString([]byte())
	// fmt.Println(sEnc)

	w.Header().Set("Content-Type", "application/json")
	prettyJSON, err4 := json.MarshalIndent(photox, "", "    ")
	if err4 != nil {
		fmt.Println(err4)
		msg := map[string]string{"message": err4.Error()}
		json.NewEncoder(w).Encode(msg)
		return
	}

	if string(prettyJSON) == "null" {
		xid := params["id"]
		msg := map[string]string{"message": "Photo ID No. " + xid + " does not exists."}
		json.NewEncoder(w).Encode(msg)
		return

	}
	w.Write(prettyJSON)
}
