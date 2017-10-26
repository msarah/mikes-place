package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Player struct holds basic player information.
type Player struct {
	Name    string
	PwHash  []byte
	Coaster string
	Wins    int
	Losses  int
	Admin   bool
	//PlayerGames []PlayerGame
}

//Session is a struct that defines a session cookie
type Session struct {
	ID       string
	Username string
}

//NewPlayer returns a pointer to a newly constructed player
func NewPlayer(Name, Coaster string, PwHash []byte, Admin bool) *Player {
	return &Player{Name, PwHash, Coaster, 0, 0, Admin}
}

//HashPassword wraps the bcrypt 'GenerateFromPassword' function and returns a byte slice of the hashed string
func HashPassword(Password string) []byte {
	bs, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	}
	return bs
}

//CreateCookie creates a session cookie for user
func CreateCookie(w http.ResponseWriter) *http.Cookie {
	id := uuid.NewV4()
	c := &http.Cookie{
		Name:     "session",
		Value:    id.String(),
		HttpOnly: true,
		MaxAge:   10080,
	}
	http.SetCookie(w, c)
	return c
}

//AlreadyLoggedIn checks if a user is logged in already
func AlreadyLoggedIn(r *http.Request) bool {
	if _, err := r.Cookie("session"); err == http.ErrNoCookie {
		return false
	}
	return true
}

//GetCookie gets the session cookie
func GetCookie(r *http.Request) *http.Cookie {
	c, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
	}
	return c
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

/*GetPlayer retrieves a player from the database
It uses the sessionID from the cookie to get the username, then takes the username and
gets the player information from the user collection*/
func (pc *PlayerController) GetPlayer(sessionID string) Player {
	var resultMap = bson.M{}
	var p = Player{}
	err := pc.SessionsCollection().Find(bson.M{"id": sessionID}).Select(bson.M{"username": 1}).One(&resultMap)
	if err != nil {
		log.Println("Could not get player: ", err)
		return p
	}

	u := resultMap["username"]

	err = pc.PlayersCollection().Find(bson.M{"name": u}).One(&p)

	return p

}

/*PlayerExist checks the database with the given username
to see if the player exists in the players collection - Returns true or false*/
func (pc *PlayerController) PlayerExist(username string) bool {
	res, err := pc.PlayersCollection().Find(bson.M{"name": username}).Count()
	if err != nil {
		log.Fatalln(err)
	}

	switch res {
	case 0:
		return false
	}
	return true
}

//GetPasswordHash retrieves the hashed password from the database associated with given username
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

/*InsertSession inserts a cookie session into the sessions collection in mongo
The cookie ID - random session ID given for each new cookie - is stored along with the username*/
func (pc *PlayerController) InsertSession(s Session) {
	err := pc.SessionsCollection().Insert(s)
	if err != nil {
		log.Fatalln(err)
	}
}

//RemoveSession removes a session from the database
func (pc *PlayerController) RemoveSession(sessionID string) {
	if err := pc.SessionsCollection().Remove(bson.M{"id": sessionID}); err != nil {
		log.Println(err)
	}

}

//RemovePlayer removes a player from the database with given name
func (pc *PlayerController) RemovePlayer(Name string) {
	if err := pc.PlayersCollection().Remove(bson.M{"name": Name}); err != nil {
		log.Println(err)
	}
}
