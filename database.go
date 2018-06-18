package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

// User entry in database
type User struct {
	Username  string
	Timestamp time.Time
	Clean     bson.Binary
	Dark      bson.Binary
}

var session *mgo.Session
var collection *mgo.Collection

// DialConnection creates sesson with mongo
func DialConnection(url string) {

	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println("can't connect to database")
	}
	collection = session.DB("flair").C("github")
}

// GetUser find the user and returns it
func GetUser(username string) User {

	result := User{}
	collection.Find(bson.M{"username": username}).One(&result)

	return result
}

// GetAllUsers returns all rows from DB
func GetAllUsers() []User {

	result := []User{}

	iter := collection.Find(nil).Iter()

	iter.All(&result)

	return result
}

// InsertUser insert the entity into DB
func InsertUser(user User) {

	collection.Insert(user)
}

// UpdateUser updated the entity
func UpdateUser(user User) {

	query := bson.M{"username": user.Username}
	change := bson.M{"$set": bson.M{
		"timestamp": user.Timestamp,
		"clean":     user.Clean,
		"dark":      user.Dark}}

	collection.Update(query, change)
}

// RemoveUser removes the entity from DB
func RemoveUser(user User) {
	collection.Remove(bson.M{"username": user.Username})
}
