package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/reynald/goapi/middleware/photos"
	"github.com/reynald/goapi/middleware/users"
)

//cors
func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/plm/cors", Cors)

	router.HandleFunc("/api/users/register", users.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/login", users.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/users/activate/{id}", users.ActivateUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users/forgotpwd/{email}", users.ForgotPassword).Methods("POST")
	router.HandleFunc("/api/users/sendtoken/{token}", users.SendToken).Methods("GET")
	router.HandleFunc("/api/users/change", users.UpdatePassword).Methods("POST")
	router.HandleFunc("/api/users/profile", users.Profile).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/photo/list", photos.List).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/photo/edit/{id}", photos.Edit).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/photo/add", photos.Add).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/photo/update/{id}", photos.Update).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/photo/delete/{id}", photos.Delete).Methods("DELETE", "OPTIONS")

	router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/")))
	log.Fatal(http.ListenAndServe(":5100", router))
}
