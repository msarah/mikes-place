package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func removeDirListHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			err := tpl.ExecuteTemplate(w, "notFound.gohtml", nil)
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
		h.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logged In: ", AlreadyLoggedIn(r))
	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatalln(err)
	}

}

func home(w http.ResponseWriter, r *http.Request) {

	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c := GetCookie(r)
	p := pc.GetPlayer(c.Value)

	if err := tpl.ExecuteTemplate(w, "home.gohtml", p); err != nil {
		log.Fatalln(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	if AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
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

func logout(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("logout func running")

	c := GetCookie(r)
	pc.RemoveSession(c.Value)
	c.MaxAge = -1
	fmt.Println("Cookie: ", c)
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
