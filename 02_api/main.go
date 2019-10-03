package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"encoding/json"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Person is ...
type Person struct {
	Fname   string `json:fname`
	Lname   string `json:lname`
	Email   string `json:email`
	Address string `json:address`
	Pword   string `json:pword`
	Id      int64  `json:id`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:password@/people")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/posts", getPerson).Methods("GET")
	router.HandleFunc("/posts", createPerson).Methods("POST")

	http.ListenAndServe(":8080", router)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []Person
	// we use “db.Query()”, this should not be used with actions that doesn’t return any rows.
	result, err := db.Query("SELECT fname, lname FROM person")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()
	// to get one record at a time from result
	for result.Next() {
		var person Person
		err := result.Scan(&person.Fname, &person.Lname)
		if err != nil {
			panic(err.Error())
		}
		people = append(people, person)
	}
	//  we encode our posts to JSON and send send them away.
	json.NewEncoder(w).Encode(people)
}
func createPerson(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	var person Person
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(body, &person)
	fmt.Println(person)
	fname := person.Fname
	lname := person.Lname
	email := person.Email
	address := person.Address
	pword := person.Pword
	id := person.Id

	stmt, err := db.Prepare("INSERT INTO person(fname,lname,email,address,pword,id) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(fname, lname, email, address, pword, id)
	if err != nil {
		panic(err.Error())
	}
	text := "Status"
	out, err := json.Marshal(text)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(out))
}
