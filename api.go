package main

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func addUser(w http.ResponseWriter, r *http.Request) {	
	fmt.Println("Endpoint hit: addUser")
	db, err := sql.Open("sqlite3", "../Prod.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	insert into User (time_created)
	values (datetime());
	`
	result, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "User added: %d", id)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/addUser", addUser)
	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}
