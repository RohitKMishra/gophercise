package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	// Underscore is used because github.com/go-sql-driver/mysql is a third party driver, so none of its exported name is visible to our code
	_ "github.com/go-sql-driver/mysql"
)

// Person is ...
type Person struct {
	Fname string `json:fname`
	Lname string `json:lname`
	Email string `json:email`
	Pword string `json:pword`
	Id    int    `json:id`
}

var db *sql.DB
var err error

func main() {
	// Creating a database object
	db, err := sql.Open("mysql", "root:password@/testdb")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	// Creating rote handling function or endpoints
	router.HandleFunc("/api/signin", getAll).Methods("GET")
	router.HandleFunc("/api/signin/{id}", getPerson).Methods("GET")
	router.HandleFunc("api/signup", createPerson).Methods("POST")
	router.HandleFunc("api/signin/{id}", updatePerson).Methods("PUT")
	router.HandleFunc("api/signin/{id}", deletePerson).Methods("DELETE")
	// starting the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// Slice of person to store record of all person
var people []Person

// to get record of all from table person
func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	// sending query over db object and storing respose in var result
	result, err := db.Query("SELECT fname, lname, email, pword, id FROM person")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// to fetch one record at a time from result
	for result.Next() {

		// creating a variable person to store the and then show it
		var person Person
		err := result.Scan(&person.Fname, &person.Lname, &person.Email, &person.Pword, &person.Id)
		if err != nil {
			panic(err.Error())
		}
		people = append(people, person)
	}
	// Encode json to be sent to client machine
	json.NewEncoder(w).Encode(people)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	pword := person.Pword
	id := person.Id

	stmt, err := db.Prepare("INSERT INTO person(fname, lname, email, pword, id) VALUES(?,?,?,?,?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(fname, lname, email, pword, id)
	if err != nil {
		panic(err.Error())
	}
	text := "Status, New user created"
	out, err := json.Marshal(text)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(out))
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var pid, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	stmt, err := db.Prepare("SELECT * FROM person WHERE id =?")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(pid)
	if err != nil {
		panic(err.Error())
	}
	_, err = result.RowsAffected()
	//json.NewEncoder(w).Encode(&Person)
}

// update existing person info

func updatePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var person Person
	fname := person.Fname
	lname := person.Lname
	email := person.Email
	pword := person.Pword
	id := person.Id

	stmt, err := db.Prepare("UPDATE person SET fname=?, lname=?, email=?, pword=?, id=?  WHERE id =?")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(fname, lname, email, pword, id)
	if err != nil {
		panic(err.Error())
	}

	_, err = result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var pid, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64) // 10,64

	stmt, err := db.Prepare("DELETE FROM person WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	result, err := stmt.Exec(pid)
	if err != nil {
		panic(err.Error())
	}
	_, err = result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
}
