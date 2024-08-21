package photos

import (
	"fmt"

	"github.com/reynald/goapi/middleware"
	"github.com/reynald/goapi/models"
)

func ValidateTitle(title string) int {
	db := middleware.CreateConnection()
	defer db.Close()

	sqlStatement := `SELECT photo_title FROM photos WHERE photo_title=?`
	// execute the sql statement
	rows, err2 := db.Query(sqlStatement, title)
	if err2 != nil {
		return 0
	}
	defer rows.Close()
	// // iterate over the rows
	var photox []models.TempPhotos
	for rows.Next() {
		var xphoto models.TempPhotos
		// unmarshal the row object to user
		err3 := rows.Scan(&xphoto.Photo_title)
		if err3 != nil {
			fmt.Println(err3)
		}
		photox = append(photox, xphoto)
	}
	rows.Close()
	if len(photox) == 0 {
		return 0
	}
	return 1
}
