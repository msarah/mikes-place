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
	/*
		password := "1234"
		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			log.Fatalln(err)
		}
		p = Player{1, "Sarah", bs, "Bud Light", 0, 0}
		pc.InsertPlayer(p)
	*/
}

func main() {

	http.Handle("/assets/", removeDirListHandler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets")))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/home", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
