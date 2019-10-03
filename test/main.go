package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type person struct {
	FirstName string `firstName: json`
	LastName  string `lastName: json`
	Email     string `email: json`
	Password  string `password: json`
}

var people []person

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeLink)
	log.Fatal(http.ListenAndServe(":8089", router))
}
