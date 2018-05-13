package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strconv"
)

var clean *image.RGBA
var dark *image.RGBA

// CreateFlair generates the flair fetching the stats from api
func CreateFlair(username string, theme string) (image.Image, error) {

	myimage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{250, 90}})

	for x := 0; x < 250; x++ {
		for y := 0; y < 90; y++ {
			if theme == "dark" {
				myimage.SetRGBA(x, y, dark.RGBAAt(x, y))
			} else {
				myimage.SetRGBA(x, y, clean.RGBAAt(x, y))
			}
		}
	}

	resp, err := HTTPGet("https://api.github.com/users/" + username)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	resp, err = HTTPGet(fmt.Sprint(data["avatar_url"]) + "")

	if err != nil {
		return nil, err
	}

	var avatar image.Image

	if resp.Header.Get("content-type") == "image/png" {
		avatar, err = png.Decode(resp.Body)
	} else {
		avatar, err = jpeg.Decode(resp.Body)
	}
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	avatar = resize.Resize(80, 80, avatar, resize.NearestNeighbor)

	// set avatar pic
	for x := 5; x <= 84; x++ {
		for y := 5; y <= 84; y++ {

			r, g, b, a := avatar.At(x-5, y-5).RGBA()
			if r == 0 && g == 0 && b == 0 && a == 0 {
				continue
			}
			//fmt.Print(avatar.At(x-5, y-5).RGBA())
			myimage.Set(x, y, avatar.At(x-5, y-5))
		}
	}

	//draw username
	dc := gg.NewContextForRGBA(myimage)
	dc.SetRGB255(3, 102, 214)
	//dc.SetRGB255(192, 192, 192)

	dc.SetFontFace(GetFontFace("arialbd.ttf", 13))

	dc.DrawStringAnchored(fmt.Sprint(data["login"]), 122, 12+6, 0, 0)
	if theme == "dark" {
		dc.SetRGB255(192, 192, 192)
	} else {
		dc.SetRGB255(38, 38, 38)
	}
	dc.DrawStringAnchored(fmt.Sprint(data["public_repos"]), 122, 12+27, 0, 0)
	dc.DrawStringAnchored(fmt.Sprint(data["followers"]), 122, 13+47, 0, 0)
	dc.DrawStringAnchored(fmt.Sprint(data["public_gists"]), 196, 13+47, 0, 0)
	//fork count and star count
	forks, stars, err := FetchCounts(username)
	if err != nil {
		return nil, err
	}
	dc.DrawStringAnchored(forks, 122, 13+67, 0, 0)
	dc.DrawStringAnchored(stars, 196, 13+27, 0, 0)

	return myimage, nil

	//jpeg.Encode(w,myimage,&jpeg.Options{Quality:100})
}

// FillIcon writes the icons to given image template
func FillIcon(im *image.RGBA, x1, y1 int, url string, theme string, errPointer *error) {

	var body []byte
	var err error

	box := packr.NewBox("./assets")
	body, _ = box.MustBytes(url)

	avatar, err := jpeg.Decode(bytes.NewReader(body))

	if err != nil && errPointer == nil {
		errPointer = &err
		return
	}

	avatar = resize.Resize(16, 16, avatar, resize.NearestNeighbor)
	for x := x1; x < x1+16; x++ {
		for y := y1; y < y1+16; y++ {
			r, g, b, a := avatar.At(x-x1, y-y1).RGBA()

			if uint8(r) == 255 && uint8(g) == 255 && uint8(b) == 255 && uint8(a) == 255 {
				im.Set(x, y, color.RGBA{238, 238, 238, 255})
				continue
			}

			if theme == "dark" && uint8(r) == 0 && uint8(g) == 0 && uint8(b) == 0 && uint8(a) == 255 {
				im.Set(x, y, color.RGBA{34, 34, 34, 255})
				continue
			}

			im.Set(x, y, avatar.At(x-x1, y-y1))
		}
	}
}

// PrepareTemplate is called first and it prepares the
// blank templates for clean and dark flairs
func PrepareTemplate() {

	fmt.Println("Preparing Template ...")

	//preparing the background image
	clean = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{250, 90}})
	dark = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{250, 90}})

	borderColor := color.RGBA{204, 204, 204, 255}

	//set the border
	for x := 0; x < 250; x++ {
		clean.Set(x, 0, borderColor)
		clean.Set(x, 89, borderColor)
		dark.Set(x, 0, borderColor)
		dark.Set(x, 89, borderColor)
	}
	for y := 0; y < 90; y++ {
		clean.Set(0, y, borderColor)
		clean.Set(249, y, borderColor)
		dark.Set(0, y, borderColor)
		dark.Set(249, y, borderColor)
	}

	//set background
	for x := 1; x < 249; x++ {
		for y := 1; y < 89; y++ {
			clean.Set(x, y, color.RGBA{238, 238, 238, 255})
			dark.Set(x, y, color.RGBA{34, 34, 34, 255})
		}
	}

	var err error

	//set the github icon
	FillIcon(clean, 97, 7, "github.jpeg", "", &err)
	FillIcon(dark, 97, 7, "github_dark.jpeg", "dark", &err)

	//set repo icon
	FillIcon(clean, 97, 28, "repo.jpeg", "", &err)
	FillIcon(dark, 97, 28, "repo_dark.jpeg", "dark", &err)

	//set followers icon
	FillIcon(clean, 97, 48, "people.jpeg", "", &err)
	FillIcon(dark, 97, 48, "people_dark.jpeg", "dark", &err)

	//set fork icon
	FillIcon(clean, 97, 68, "fork.jpeg", "", &err)
	FillIcon(dark, 97, 68, "fork_dark.jpeg", "dark", &err)

	//set gist icon
	FillIcon(clean, 173, 48, "gist.jpeg", "", &err)
	FillIcon(dark, 173, 48, "gist_dark.jpeg", "dark", &err)

	//set star icon
	FillIcon(clean, 173, 28, "star.jpeg", "", &err)
	FillIcon(dark, 173, 28, "star_dark.jpeg", "dark", &err)

	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}

// FetchCounts return the total fork and star count of every
// repo of the user.
func FetchCounts(username string) (string, string, error) {

	resp, err := HTTPGet("https://api.github.com/users/" + username + "/repos")

	if err != nil {
		return "", "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var data1 []map[string]interface{}
	json.Unmarshal(body, &data1)

	forkCount, starCount := 0, 0
	for v := range data1 {
		temp, _ := strconv.Atoi(fmt.Sprint(data1[v]["forks_count"]))
		forkCount = forkCount + temp

		temp, _ = strconv.Atoi(fmt.Sprint(data1[v]["stargazers_count"]))
		starCount = starCount + temp
	}
	return strconv.Itoa(forkCount), strconv.Itoa(starCount), nil
}

// HTTPGet returns the json from url
func HTTPGet(url string) (*http.Response, error) {
	var resp *http.Response

	box := packr.NewBox("./secrets")

	body, err := box.MustString("access_token.txt")

	if err == nil {
		url = url + "?access_token=" + body[:len(body)-1]
	}

	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetFontFace return an instance of font.Face using the arial
// font included in the assests folder
func GetFontFace(filename string, points float64) font.Face {

	box := packr.NewBox("./assets")

	fontBytes, err := box.MustBytes(filename)

	if err != nil {
		panic(err)
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face
}
