package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Player struct holds basic player information.
type Player struct {
	ID      int16
	Name    string
	PwHash  []byte
	Coaster string
	Wins    int
	Losses  int
}

//PlayerController is a srtuct through which the database is safely accessed
type PlayerController struct {
	session *mgo.Session
}

//NewPlayerController returns a pointer to a PlayerController for use
func NewPlayerController(s *mgo.Session) *PlayerController {
	return &PlayerController{s}
}

//PlayersCollection returns a pointer to the players mongo collection
func (pc *PlayerController) PlayersCollection() *mgo.Collection {
	return pc.session.DB(db).C("players")
}

//SessionsCollection returns a pointer to the sessions mongo collection
func (pc *PlayerController) SessionsCollection() *mgo.Collection {
	return pc.session.DB(db).C("sessions")
}

/*Exist checks the database to see if the given data is available.
Returns true or false*/
func (pc *PlayerController) Exist(key, value string) bool {
	res, err := pc.PlayersCollection().Find(bson.M{key: value}).Count()
	if err != nil {
		log.Fatalln(err)
	}

	switch res {
	case 0:
		return false
	}
	return true
}

//GetPasswordHash retrieves the hashed password associated with given username
func (pc *PlayerController) GetPasswordHash(username string) []byte {

	//Scan result from mongo search into resultMap
	var resultMap = bson.M{}

	//Finds the document containing given username
	//Returns the password hash associated with username in database
	err := pc.PlayersCollection().Find(bson.M{"name": username}).Select(bson.M{"pwhash": 1}).One(&resultMap)
	if err != nil {
		log.Println("No password hash with given username", err)
	}

	//Retrieve password hash from result map using key
	bs := resultMap["pwhash"].([]byte)

	return bs
}

//InsertPlayer inserts a player into the player collection in mongo
func (pc *PlayerController) InsertPlayer(p Player) {
	if err := pc.PlayersCollection().Insert(p); err != nil {
		log.Fatalln(err)
	}
}
