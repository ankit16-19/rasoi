package dbConnection

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

// DAO :
type DAO struct{}

// Db :
var Db *mgo.Database

// Connect :
func (c *DAO) Connect() {
	session, err := mgo.Dial("172.16.1.213")
	if err != nil {
		log.Fatal(err)
	}
	Db = session.DB("rasoi")

}
