package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, r *http.Request) {

	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatalln(err)
	}

}

func login(w http.ResponseWriter, r *http.Request) {

	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {
		u := r.FormValue("username")
		p := r.FormValue("password")

		if pc.PlayerExist(u) {
			hp := pc.GetPasswordHash(u)
			err := bcrypt.CompareHashAndPassword(hp, []byte(p))
			if err != nil {
				log.Fatalln(err)
			}
		}

		c := CreateCookie(w)
		s := Session{c.Value, u}
		pc.InsertSession(s)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	c := GetCookie(r)
	p := pc.GetPlayer(c.Value)

	if err := tpl.ExecuteTemplate(w, "home.gohtml", p); err != nil {
		log.Fatalln(err)
	}
}
