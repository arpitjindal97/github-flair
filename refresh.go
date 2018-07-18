package main

import (
	"bytes"
	"image/png"
	"log"
	"time"
)

// RefreshImages is called after an interval to
// - Remove the entries of non-active users
// - Refresh the flairs of active users
func RefreshImages() {

	result := GetAllUsers()

	log.Println("Starting refresh ...")
	for _, user := range result {

		hours := time.Now().Sub(user.Timestamp).Hours()

		log.Println("Refreshing " + user.Username)

		// remove user non-active for 5 days
		if hours > 120 {

			log.Println("Deleted")
			RemoveUser(user)

		} else {
			log.Println("Updated")
			image, _ := CreateFlair(user.Username, "clean")
			buffer := new(bytes.Buffer)
			png.Encode(buffer, image)
			user.Clean.Data = buffer.Bytes()

			image, _ = CreateFlair(user.Username, "dark")
			buffer = new(bytes.Buffer)
			png.Encode(buffer, image)
			user.Dark.Data = buffer.Bytes()

			UpdateUser(user)

		}
	}
	log.Println("Refresh completed")

}
