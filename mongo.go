package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

const db = "mikes-place"

//GetSession retrieves a mongo session
func GetSession() *mgo.Session {

	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		log.Fatalln(err)
	}
	return session
}
