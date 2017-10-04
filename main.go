package main

import (
	"html/template"
	"log"
	"net/http"
)

var (
	tpl *template.Template
	pc  *PlayerController
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	pc = NewPlayerController(GetSession())
}

func main() {

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
