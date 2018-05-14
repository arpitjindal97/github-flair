package main

import (
	"bytes"
	"image/png"
	"time"
)

// RefreshImages is called after an interval to
// - Remove the entries of non-active users
// - Refresh the flairs of active users
func RefreshImages() {

	result := GetAllUsers()

	for _, user := range result {

		hours := time.Now().Sub(user.Timestamp).Hours()

		// remove user non-active for 5 days
		if hours > 120 {

			RemoveUser(user)

		} else {

			image, _ := CreateFlair(user.Username, "clean")
			buffer := new(bytes.Buffer)
			png.Encode(buffer, image)
			user.Clean.Data = buffer.Bytes()

			image, _ = CreateFlair(user.Username, "dark")
			buffer = new(bytes.Buffer)
			png.Encode(buffer, image)
			user.Dark.Data = buffer.Bytes()

			user.Timestamp = time.Now()
			UpdateUser(user)

		}
	}

}
