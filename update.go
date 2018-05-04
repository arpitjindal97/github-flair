package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"image/png"
	"os"
	"time"
)

func RefreshImages() {

	session, _ := mgo.Dial("127.0.0.1")
	defer session.Close()

	collection := session.DB("flair").C("github")

	result := []User{}

	iter := collection.Find(nil).Iter()

	iter.All(&result)

	PrepareTemplate()

	for _, user := range result {

		os.Remove("flairs/" + user.Username + ".png.clean")
		os.Remove("flairs/" + user.Username + ".png.dark")

		hours := time.Now().Sub(user.Timestamp).Hours()
		if hours > 120 {

			collection.Remove(bson.M{"username": user.Username})

		} else {

			file, _ := os.Create("flairs/" + user.Username + ".png.clean")
			png.Encode(file, CreateFlair(user.Username, "clean"))
			defer file.Close()

			file1, _ := os.Create("flairs/" + user.Username + ".png.dark")
			png.Encode(file1, CreateFlair(user.Username, "dark"))

			defer file1.Close()
		}
		fmt.Println(user)
	}

}
