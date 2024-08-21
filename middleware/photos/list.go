// http://localhost:5100/photo/list

package photos

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

//LIST PHOTO
func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//AUTHORIZATION
	_, err := middleware.ExtractTokenID(r)
	if err != nil {
		// middleware.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		msg := map[string]string{"message": "UnAuthorized Access."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	db := middleware.CreateConnection()
	defer db.Close()
	var photox []models.TempPhotos

	// create the select sql query
	sqlStatement := `SELECT id, photo_title, photo_image FROM photos ORDER BY id`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
		msg := map[string]string{"message": "Unable to execute the query."}
		json.NewEncoder(w).Encode(msg)
		return
	}
	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var xphoto models.TempPhotos

		// unmarshal the row object to user
		err = rows.Scan(&xphoto.ID, &xphoto.Photo_title, &xphoto.Photo_image)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		photox = append(photox, xphoto)
	}
	w.Header().Set("Content-Type", "application/json")
	prettyJSON, err := json.MarshalIndent(photox, "", "    ")
	if err != nil {
		msg := map[string]string{"message": err.Error()}
		json.NewEncoder(w).Encode(msg)
	}
	w.Write(prettyJSON)
}
