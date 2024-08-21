package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var dbVideos *sql.DB
var dbase2 = "rey:rey@tcp(127.0.0.1:3306)/godb"

//Videos Model
type VTag struct {
	ID          string `json:"id"`
	Video_title string `json:"video_title"`
	Video_file  string `json:"video_file"`
}

type Videos struct {
	ID         string `json:"id"`
	Videotitle string `json:"videotitle"`
	Videofile  string `json:"videofile"`
}

var videos []Videos

//index
func getVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)

}
func getVideo(w http.ResponseWriter, r *http.Request) {

	//show
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	for _, item := range videos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

//Create
func createVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newVideo Videos
	json.NewDecoder(r.Body).Decode(&newVideo)
	newVideo.ID = strconv.Itoa(len(videos) + 1)
	videos = append(videos, newVideo)
	json.NewEncoder(w).Encode(newVideo)

	sql := `INSERT INTO videos (id,video_title,video_image) VALUES ('2','Graduation Day','002.jpg')`
	dbVideos.Exec(sql)

}

//Update
func updateVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Println("video update test...")
}

//Delete
func deleteVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("video delelete test...")
}

func handleVideosRequest() {
	dbVideos, err2 := sql.Open("mysql", dbase2)
	if err2 != nil {
		log.Println(err2.Error())
	}
	defer dbVideos.Close()
	results, err := dbVideos.Query("SELECT id, video_title, video_file FROM videos")
	if err != nil {
		log.Println(err.Error())
	}

	for results.Next() {
		var vtag VTag
		err = results.Scan(&vtag.ID, &vtag.Video_title, &vtag.Video_file)
		if err != nil {
			log.Println(err.Error())
		}
		videos = append(videos, Videos{ID: vtag.ID, Videotitle: vtag.Video_title, Videofile: vtag.Video_file})
	}
}
