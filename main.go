package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

var db *sql.DB

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	connectionString := "admin:13731373@tcp(database-1.ctgevmgdpdm8.us-east-1.rds.amazonaws.com:3306)/Person"
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/", rootSite).Methods("GET")
	r.HandleFunc("/people", GetPeople).Methods("GET")
	r.HandleFunc("/people/{id:[0-9]+}", GetPerson).Methods("GET")
	r.HandleFunc("/people", CreatePerson).Methods("POST")
	r.HandleFunc("/people/{id:[0-9]+}", UpdatePerson).Methods("PUT")
	r.HandleFunc("/people/{id:[0-9]+}", DeletePerson).Methods("DELETE")
	port := 9090
	fmt.Printf("Server started on port %d..", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))

}
func rootSite(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
func GetPeople(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to retrieve all people from the database
	rows, err := db.Query("select id,name,age from information ")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var people []Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.ID, &p.Name, &p.Age)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		}
		people = append(people, p)
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to retrieve a single person by ID from the database
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	var p Person
	err = db.QueryRow("SELECT id, name, age FROM information WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}
		log.Fatal(err)
		http.Error(w, "fail to fetvh data!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/jason")
	json.NewEncoder(w).Encode(p)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to create a new person in the database
	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Fail to create person", http.StatusInternalServerError)
		return
	}
	reuslt, err := db.Exec("INSERT INTO information (name, age,id) VALUES (?, ?, ?)", p.Name, p.Age, p.ID)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}
	lastInsertID, err := reuslt.LastInsertId()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to get last insert ID", http.StatusInternalServerError)
		return

	}

	p.ID = int(lastInsertID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to update a person by ID in the database
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p Person
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Update the person in the database
	_, err = db.Exec("UPDATE information SET name = ?, age = ? WHERE id = ?", p.Name, p.Age, id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to delete a person by ID from the database
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete the person from the database
	_, err = db.Exec("DELETE FROM information WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
