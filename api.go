package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func open() (*sql.DB, error) {
	return sql.Open("sqlite3", "../Prod.db")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func addUser(w http.ResponseWriter, r *http.Request) {
	db, err := open()
	if err != nil {
		fmt.Printf("Error in addUser, open(): %s\n", err)
		return
	}
	defer db.Close()

	stmt := `
	insert into User (time_created)
	values (datetime());
	`
	result, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Error in addUser, db.Exec(): %s\n", err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Error in addUser, result.LastInsertId(): %s\n", err)
		return
	}
	fmt.Fprintf(w, "User added: %d", id)
}

func addAnalyticsEvent(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error in addAnalyticsEvent, Decode: %s\n", err)
		return
	}
	userID := data["userID"]
	eventType := data["type"]

	db, err := open()
	if err != nil {
		fmt.Printf("Error in addAnalyticsEvent, open(): %s\n", err)
		return
	}
	defer db.Close()

	stmt := `
	insert into AnalyticsEvent (user_id, type, time)
	values (?, ?, datetime());
	`
	_, err = db.Exec(stmt, userID, eventType)
	if err != nil {
		fmt.Printf("Error in addAnalyticsEvent, db.Exec(): %s\n", err)
		return
	}
	fmt.Fprintf(w, `{"success": true}`)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Endpoint hit: %s\n", pattern)
		enableCors(&w)
		handler(w, r)
	})
}

func handleRequests() {
	handleFunc("/", homePage)
	handleFunc("/addUser", addUser)
	handleFunc("/addAnalyticsEvent", addAnalyticsEvent)
	fmt.Println("Listening at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	handleRequests()
}
