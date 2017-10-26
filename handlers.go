package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//ParseBool returns true for a string of '1' and false for anything else
func ParseBool(str string) bool {
	if str == "1" {
		return true
	}
	return false
}

//RedirectHome wraps the http Redirect function
func RedirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

//RedirectIndex wraps the http Redirect function
func RedirectIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index", http.StatusSeeOther)
}

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
	if AlreadyLoggedIn(r) {
		RedirectHome(w, r)
	}

	if err := tpl.ExecuteTemplate(w, "index.gohtml", nil); err != nil {
		log.Fatalln(err)
	}

}

func home(w http.ResponseWriter, r *http.Request) {

	if !AlreadyLoggedIn(r) {
		RedirectIndex(w, r)
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
		RedirectHome(w, r)
		return
	}

	if r.Method == http.MethodPost {
		u := r.FormValue("username")
		p := r.FormValue("password")

		if pc.PlayerExist(u) {
			hp := pc.GetPasswordHash(u)
			err := bcrypt.CompareHashAndPassword(hp, []byte(p))
			if err != nil {
				log.Println("Wrong Password!", err)
				RedirectHome(w, r)
				return
			}
			c := CreateCookie(w)
			s := Session{c.Value, u}
			pc.InsertSession(s)
			RedirectHome(w, r)
			return
		}
		fmt.Println("sorry cannot find player!")
		RedirectHome(w, r)
		return
	}
	fmt.Println("sorry cannot find player!")
	RedirectHome(w, r)

}

func logout(w http.ResponseWriter, r *http.Request) {
	if !AlreadyLoggedIn(r) {
		RedirectIndex(w, r)
		return
	}

	c := GetCookie(r)
	pc.RemoveSession(c.Value)
	c.MaxAge = -1
	http.SetCookie(w, c)
	RedirectIndex(w, r)
}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	var p *Player

	if r.Method == http.MethodPost {
		n := r.FormValue("name")
		c := r.FormValue("coaster")
		pw := HashPassword(r.FormValue("password"))
		a := ParseBool(r.FormValue("admin"))
		p = NewPlayer(n, c, pw, a)
		pc.InsertPlayer(*p)
	}

	if err := tpl.ExecuteTemplate(w, "addPlayer.gohtml", p); err != nil {
		log.Println(err)
	}
}

func removePlayer(w http.ResponseWriter, r *http.Request) {
	var n string
	if r.Method == http.MethodPost {
		n = r.FormValue("name")
		pc.RemovePlayer(n)
	}
	if err := tpl.ExecuteTemplate(w, "removePlayer.gohtml", n); err != nil {
		log.Println(err)
	}
}
