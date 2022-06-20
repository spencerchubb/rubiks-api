package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
)

func logEndpoint(endpoint string) {
	fmt.Printf("Endpoint hit: %s\n", endpoint)
}

func open() *sql.DB {
	return sql.Open("sqlite3", "../Prod.db")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	logEndpoint("homePage")
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func addUser(w http.ResponseWriter, r *http.Request) {
	logEndpoint("addUser")

	db, err := open()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	stmt := `
	insert into User (time_created)
	values (datetime());
	`
	result, err := db.Exec(stmt)
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

func addAnalyticsEvent(w http.ResponseWriter, r *http.Request) {
	logEndpoint("addAnalyticsEvent")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Printf("Could not unmarshal json: %s\n", err)
		return
	}
	userID := data["userID"]
	eventType := data["type"]

	db, err := open()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	stmt := `
	insert into AnalyticsEvent (user_id, type, time)
	values (?, ?, datetime());
	`
	_, err = db.Exec(stmt, userID, eventType)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, `{"success": true}`)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/addUser", addUser)
	http.HandleFunc("/addAnalyticsEvent", addAnalyticsEvent)
	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}
