package main

import (
	"BE_WE_SAVING/internal/app"
	"log"
	"net/http"
	"os"
)

func main() {
	s := app.NewServer()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Router))
}
