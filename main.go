package main

import (
	"bytes"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mileusna/crontab"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	PrepareTemplate()

	_, err := ioutil.ReadDir("flairs")
	if err != nil {
		os.Mkdir("flairs", 0755)
	}

	ctab := crontab.New()
	ctab.AddJob("10 23 * * *", RefreshImages)

	http.HandleFunc("/github/", Flair)

	http.ListenAndServeTLS(":443", "crt-bundle.pem",
		"ssl-private.key", nil)

	// http.ListenAndServe(":8080", nil)

}

func Flair(w http.ResponseWriter, r *http.Request) {

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

type User struct {
	Username  string
	Timestamp time.Time
}

func ExistsInDatabase(username string) bool {

	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()

	collection := session.DB("flair").C("github")

	result := User{}

	collection.Find(bson.M{"username": username}).One(&result)

	if result.Username == "" {
		return false
	}
	return true
}

func UpdateDatabase(username string) {

	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()

	collection := session.DB("flair").C("github")

	query := bson.M{"username": username}
	change := bson.M{"$set": bson.M{"timestamp": time.Now()}}

	collection.Update(query, change)

}

func InsertDatabase(username string) {

	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()

	collection := session.DB("flair").C("github")

	user := &User{Username: username, Timestamp: time.Now()}

	collection.Insert(user)

	file, _ := os.Create("flairs/" + username + ".png.clean")

	png.Encode(file, CreateFlair(username, "clean"))

	file.Close()

	file, _ = os.Create("flairs/" + username + ".png.dark")

	png.Encode(file, CreateFlair(username, "dark"))

	defer file.Close()

}

func GetFromFolder(username string, theme string) image.Image {

	file, _ := os.Open("flairs/" + username + ".png." + theme)

	defer file.Close()

	image, _ := png.Decode(file)

	return image
}
