// http://localhost:5100/photo/delete/90

package photos

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/reynald/goapi/middleware"
)

//DELETE PHOTO
func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//AUTHORIZATION
	_, err := middleware.ExtractTokenID(r)
	if err != nil {
		msg := map[string]string{"message": "UnAuthorized Access."}
		json.NewEncoder(w).Encode(msg)
		return
	}

	params := mux.Vars(r)
	idno, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	db := middleware.CreateConnection()
	defer db.Close()

	//INSERT
	// insertcontact($1, $2, $3, $4, $5)`
	res, errs2 := db.Exec("DELETE FROM photos WHERE id=?", idno)
	if errs2 != nil {
		msg := map[string]string{"message": errs2.Error()}
		json.NewEncoder(w).Encode(msg)
		return
	}
	xid := params["id"]
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected > 0 {
		msg := map[string]string{"message": "Photo ID No." + xid + " Deleted Successfully."}
		json.NewEncoder(w).Encode(msg)
	} else {
		msg := map[string]string{"message": "Photo ID No. " + xid + " does not exists."}
		json.NewEncoder(w).Encode(msg)
	}

}
