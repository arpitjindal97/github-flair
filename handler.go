package main

import (
	"bytes"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Flair handles /github/ URL requests
// The URL should be  in this format
// github/username.png It extracts the username and finds it in
// database, If it exists then image is
// provided from folder else it os created
// and put in folder and returned
func Flair(w http.ResponseWriter, r *http.Request) {

	// extracts username from url
	username := strings.SplitN(r.URL.Path[1:], "/", 2)[1]
	username = username[:len(username)-4]

	theme := r.URL.Query().Get("theme")
	if theme == "" {
		theme = "clean"
	}

	log.Println(username)

	var myimage image.Image

	if ExistsInDatabase(username) == true {
		myimage = GetFromFolder(username, theme)
		UpdateDatabase(username)
	} else {
		myimage = CreateFlair(username, theme)
		InsertDatabase(username)
	}

	w.Header().Set("Content-Type", "image/png")
	buffer := new(bytes.Buffer)
	png.Encode(buffer, myimage)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	w.Write(buffer.Bytes())
}

// User is an entry of a single user in Database
type User struct {
	Username  string
	Timestamp time.Time // last time opened
}

// ExistsInDatabase tells if user entry exists
// in Database or not
func ExistsInDatabase(username string) bool {

	session, _ := mgo.Dial("mongo")
	defer session.Close()

	collection := session.DB("flair").C("github")

	result := User{}

	collection.Find(bson.M{"username": username}).One(&result)

	if result.Username == "" {
		return false
	}
	return true
}

// UpdateDatabase updates the existing entry with
// updated timestamp
func UpdateDatabase(username string) {

	session, _ := mgo.Dial("mongo")
	defer session.Close()

	collection := session.DB("flair").C("github")

	query := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"timestamp": time.Now()}}

	collection.Update(query, change)

}

// InsertDatabase inserts a new entry of user in
// Database with current timestamp
func InsertDatabase(username string) {

	session, _ := mgo.Dial("mongo")
	defer session.Close()

	collection := session.DB("flair").C("github")

	user := &User{Username: username, Timestamp: time.Now()}

	collection.Insert(user)

	PutInFolder(username)
}

// PutInFolder generates the image and puts it
// in the folder
func PutInFolder(username string) {

	file, _ := os.Create("/data/flair-images/" + username + ".png.clean")

	png.Encode(file, CreateFlair(username, "clean"))

	file.Close()

	file, _ = os.Create("/data/flair-images/" + username + ".png.dark")

	png.Encode(file, CreateFlair(username, "dark"))

	defer file.Close()
}

// GetFromFolder fetches the image from folder
// it will put the image also if not found
func GetFromFolder(username string, theme string) image.Image {

	file, err := os.Open("/data/flair-images/" + username + ".png." + theme)

	defer file.Close()

	if err != nil {
		log.Println("Error! Image not found in folder, recreating ..")
		PutInFolder(username)
		return CreateFlair(username, theme)
	}
	image, _ := png.Decode(file)

	return image
}
