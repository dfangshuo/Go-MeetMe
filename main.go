package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"       // log errors and stuff
	"math/rand" // specific to tutorial: generate random num
	"net/http"  // work with http
	"strconv"   // str converter
)

// MODELS
// Book Struct
type Book struct {
	Id     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {

	/*
		Let's take a step back and try to see what's happening here:
		your r is your REQUEST
		your w is your RESPONSE
		This is saying, firstly, SET (in the header of your response) that you're sending back at JSON
		ok. This is so you won't be serving content as txt
		THEN, you want to create a new encder that WRITES TO w (your response)
		And specifically, you want to write your BOOKS to your response
	*/

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	// .NewEncoder returns a NewEncoder that writes to w
	// Encode writes the JSON encoding of v (books) into the stream (the new Encoder returned )
}

// get 1 book
func getBook(w http.ResponseWriter, r *http.Request) {
	//parse id from r
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get any params
	for _, book := range books {
		if book.Id == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	// empty book for book not found
	json.NewEncoder(w).Encode(&Book{Title: "ERROR: book not found!"})

}

// create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook Book

	/*
		NewDecoder creates a NewDecoder that READS from the source (r.Body, in this case)
		Decode reads the next JSON-encoded value from its input and
		stores it in the value pointed to by v (which in this case, is &newBook)
	*/
	_ = json.NewDecoder(r.Body).Decode(&newBook)
	newBook.Id = strconv.Itoa(rand.Intn(100000000)) // not safe for prod - overlapping Ids
	// adds the book created to the 'database'
	books = append(books, newBook)
	// returns the new book created
	json.NewEncoder(w).Encode(newBook)

}

// update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get any params
	for i, book := range books {
		if book.Id == params["id"] {
			books = append(books[:i], books[i+1:]...)
			var newBook Book
			_ = json.NewDecoder(r.Body).Decode(&newBook)
			newBook.Id = params["id"]
			books = append(books, newBook)
			json.NewEncoder(w).Encode(newBook)
			return
		}
	}
	// remove that book
	json.NewEncoder(w).Encode(books)
}

// delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get any params
	for i, book := range books {
		if book.Id == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	// remove that book
	json.NewEncoder(w).Encode(books)
}

func main() {
	// init router
	r := mux.NewRouter()

	// Mock Data
	// TODO: implement database
	books = append(books, Book{
		Id:    "1",
		Isbn:  "94704",
		Title: "Book One",
		Author: &Author{
			Firstname: "John",
			Lastname:  "W",
		},
	})

	books = append(books, Book{
		Id:    "2",
		Isbn:  "21421",
		Title: "Book Two",
		Author: &Author{
			Firstname: "Steve",
			Lastname:  "Smith",
		},
	})

	// Route handlers/ endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
