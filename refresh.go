package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"image/png"
	"os"
	"time"
)

// RefreshImages is called after an interval to
// - Remove the entries of non-active users
// - Refreshe the flairs of active users
func RefreshImages() {

	session, _ := mgo.Dial("localhost")
	defer session.Close()

	collection := session.DB("flair").C("github")

	result := []User{}

	iter := collection.Find(nil).Iter()

	iter.All(&result)

	for _, user := range result {

		os.Remove("data-db/flair-images/" + user.Username + ".png.clean")
		os.Remove("data-db/flair-images/" + user.Username + ".png.dark")

		hours := time.Now().Sub(user.Timestamp).Hours()

		// remove user non-active for 5 days
		if hours > 120 {

			collection.Remove(bson.M{"username": user.Username})

		} else {

			file, _ := os.Create("data-db/flair-images/" + user.Username + ".png.clean")
			img, _ := CreateFlair(user.Username, "clean")
			png.Encode(file, img)
			defer file.Close()

			file1, _ := os.Create("data-db/flair-images/" + user.Username + ".png.dark")
			img, _ = CreateFlair(user.Username, "clean")
			png.Encode(file1, img)

			defer file1.Close()
		}
		fmt.Println(user)
	}

}
