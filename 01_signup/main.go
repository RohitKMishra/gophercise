package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type person struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func signup(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-type", "applicaton/json")

	json.NewDecoder(router.Body).Decode(&signup)

}

var people []person

func main() {

	router := mux.NewRouter()

	// i := info{
	// 	FirstName: "James",
	// 	LastName:  "Bond",
	// 	Email:     "jamesbond@gmail.com",
	// 	Password:  "bond007",
	// }
	router.HandleFunc("/signup", signup).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
