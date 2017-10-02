package main

import (
	"log"
	"net/http"
)

func main() {

	http.StripPrefix("/assets/", http.FileServer(http.Dir()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
