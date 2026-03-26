package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"Book-API_Golang/handlers"
	"Book-API_Golang/models"
)

func main() {
	store := models.NewStore()

	r := mux.NewRouter()

	bookHandler := &handlers.BookHandler{Store: store}
	authorHandler := &handlers.AuthorHandler{Store: store}
	categoryHandler := &handlers.CategoryHandler{Store: store}

	bookHandler.RegisterRoutes(r)
	authorHandler.RegisterRoutes(r)
	categoryHandler.RegisterRoutes(r)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
