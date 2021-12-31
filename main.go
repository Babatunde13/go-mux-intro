package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type BookResponse struct {
	Message string `json:"message"`
	Data Book `json:"data"`
	Status bool `json:"status"`
}

type BooksResponse struct {
	Message string `json:"message"`
	Data []Book `json:"data"`
	Status bool `json:"status"`
}

type NotFound struct {
	Message string `json:"message"`
	Status bool `json:"status"`
}

var books = []Book{
	{ID: "1", Isbn: "123", Title: "Book One", Author: Author{Firstname: "John", Lastname: "Doe"}},
	{ID: "2", Isbn: "456", Title: "Book Two", Author: Author{Firstname: "Steve", Lastname: "Smith"}},
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BooksResponse{Message: "Books fetched successfully", Data: books, Status: true})
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	for _, book := range books {
		if book.ID == id {
			w.WriteHeader(http.StatusOK)
			resp := BookResponse{Message: "Book fetched successfully", Data: book, Status: true}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(NotFound{Message: "Book not found", Status: false})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(BookResponse{Message: "Book created successfully", Data: book, Status: true})
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	for index, book := range books {
		if book.ID == id {
			var updateBook Book
			_ = json.NewDecoder(r.Body).Decode(&updateBook)
			if updateBook.Isbn != "" && updateBook.Isbn != book.Isbn {
				book.Isbn = updateBook.Isbn
			}
			if updateBook.Title != "" && updateBook.Title != book.Title {
				book.Title = updateBook.Title
			}
			if updateBook.Author.Firstname != "" && updateBook.Author.Firstname != book.Author.Firstname {
				book.Author.Firstname = updateBook.Author.Firstname
			}
			if updateBook.Author.Lastname != "" && updateBook.Author.Lastname != book.Author.Lastname {
				book.Author.Lastname = updateBook.Author.Lastname
			}
			books = append(append(books[:index], book), books[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(BookResponse{Message: "Book updated successfully", Data: book, Status: true})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(NotFound{Message: "Book not found", Status: false})
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	for index, book := range books {
		if book.ID == id {
			books = append(books[:index], books[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(BookResponse{Message: "Book deleted successfully", Data: book, Status: true})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(NotFound{Message: "Book not found", Status: false})
}


func main () {
	// initialize the router
	r := mux.NewRouter()

	// Route handlers & endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}