package main

import (
	"context"
	"encoding/json"
	"firebase.google.com/go/db"
	"github.com/dfangshuo/mmModels"
	"github.com/gorilla/mux"
	"log"      // log errors and stuff
	"net/http" // work with http
	// "strconv"  // str converter
	"github.com/kjk/betterguid"

	firebase "firebase.google.com/go"
	// "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// Init books var as a slice Book struct
var users []mmModels.MeetMeUser
var port = ":8000"
var client *db.Client
var userRef *db.Ref
var ctx context.Context
var router *mux.Router

// get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []mmModels.MeetMeUser
	if err := userRef.Get(ctx, &users); err != nil {
		log.Fatalf("getUsers error | %v\n", err)
	}
	json.NewEncoder(w).Encode(users)
}

// HELPER function to get 1 user by ID
func getUserByID(id string) mmModels.MeetMeUser {
	var user mmModels.MeetMeUser
	if err := userRef.Child(id).Get(ctx, &user); err != nil {
		log.Fatalf("getUserByID error | %v\n", err)
	}
	return user
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get any params
	user := getUserByID(params["id"])
	json.NewEncoder(w).Encode(user)
}

// create a new book
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser mmModels.MeetMeUser
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	newID := betterguid.New()

	for getUserByID(newID).ID != "" {
		newID = betterguid.New()
	}

	newUser.ID = newID

	if err := userRef.Child(newID).Set(ctx, newUser); err != nil {
		log.Fatalf("ref error, %v", err)
	}

	json.NewEncoder(w).Encode(newUser)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if err := userRef.Child(id).Delete(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// // update book
// func updateBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r) // get any params
// 	for i, book := range books {
// 		if book.Id == params["id"] {
// 			books = append(books[:i], books[i+1:]...)
// 			var newBook mmModels.Book
// 			_ = json.NewDecoder(r.Body).Decode(&newBook)
// 			newBook.Id = params["id"]
// 			books = append(books, newBook)
// 			json.NewEncoder(w).Encode(newBook)
// 			return
// 		}
// 	}
// 	// remove that book
// 	json.NewEncoder(w).Encode(books)
// }

// // delete book
// func deleteBook(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r) // get any params
// 	for i, book := range books {
// 		if book.Id == params["id"] {
// 			books = append(books[:i], books[i+1:]...)
// 			break
// 		}
// 	}
// 	// remove that book
// 	json.NewEncoder(w).Encode(books)
// }

func main() {
	ctx = context.Background()
	config := &firebase.Config{
		DatabaseURL: "https://gotest-76eff.firebaseio.com/",
	}
	opt := option.WithCredentialsFile("gotest-76eff-firebase-adminsdk-pjcks-795d67e1db.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("app initialization failure | %v", err)
	}
	client, err = app.Database(ctx)
	if err != nil {
		log.Fatalf("client initialization failure | %v", err)
	}
	userRef = client.NewRef("Users")

	router = mux.NewRouter()
	// Route handlers/ endpoints
	// router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getOneUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	// r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(port, router))
}
