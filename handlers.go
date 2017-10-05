package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatalln(err)
	}

	pc.session.Ping()
}

func login(w http.ResponseWriter, r *http.Request) {
	var u string
	var p string

	if r.Method == http.MethodPost {
		u = r.FormValue("username")
		p = r.FormValue("password")

		if pc.Exist("name", u) {
			hp := pc.GetPasswordHash(u)
			err := bcrypt.CompareHashAndPassword(hp, []byte(p))
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("passwords match!")
		}

		tpl.ExecuteTemplate(w, "home.gohtml", u)
	}

}
