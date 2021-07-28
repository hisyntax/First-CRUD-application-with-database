package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var DB *gorm.DB
var err error

type User struct {
	gorm.Model
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func initailMigration() {
	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to database")
	}
	defer DB.Close()

	DB.AutoMigrate(&User{})
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot get Users from database")
	}
	defer DB.Close()

	var users []User
	DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot get User from the database")
	}
	defer DB.Close()

	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot create new User to the database")
	}
	defer DB.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)
	DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot update User in database")
	}
	defer DB.Close()

	params := mux.Vars(r)
	var user User
	DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)
	DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DB, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot delete User from database")
	}
	defer DB.Close()

	params := mux.Vars(r)
	var user User
	DB.Delete(&user, params["id"])
	json.NewEncoder(w).Encode("This user has been successfully updated!")
}

func initializeRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func main() {
	fmt.Println("Server running at port :8000")
	initailMigration()
	initializeRouter()
}
