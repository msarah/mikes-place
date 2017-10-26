package main

import (
	"html/template"
	"log"
	"net/http"
)

var (
	tpl *template.Template
	pc  *PlayerController
	p   Player
)

func init() {
	tpl = template.Must(template.ParseGlob("assets/templates/*")) //variable that stores access to templates
	pc = NewPlayerController(GetSession())
	//pc.InsertPlayer(*NewPlayer("Mike", "Samuel Adams", 0, 0, HashPassword("mikesplace"), false))
}

func main() {

	http.Handle("/assets/", removeDirListHandler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/home", home)
	http.HandleFunc("/addPlayer", addPlayer)
	http.HandleFunc("/removePlayer", removePlayer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
