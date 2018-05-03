package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gobuffalo/packr"
	"github.com/nfnt/resize"
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

func CreateFlair(username string, theme string) image.Image {

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

	resp := HttpGet("https://api.github.com/users/" + username)

	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	resp = HttpGet(fmt.Sprint(data["avatar_url"]) + "")

	var avatar image.Image
	var err error

	if resp.Header.Get("content-type") == "image/png" {
		avatar, err = png.Decode(resp.Body)
	} else {
		avatar, err = jpeg.Decode(resp.Body)
	}
	defer resp.Body.Close()

	if err != nil {
		panic(err)
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
	dc.LoadFontFace("arialbd.ttf", 13)

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
	forks, stars := FetchCounts(username)
	dc.DrawStringAnchored(forks, 122, 13+67, 0, 0)
	dc.DrawStringAnchored(stars, 196, 13+27, 0, 0)

	return myimage

	//jpeg.Encode(w,myimage,&jpeg.Options{Quality:100})
}

func FillIcon(im *image.RGBA, x1, y1 int, url string, theme string) {

	// resp := HttpGet(url)
	// avatar, _ := jpeg.Decode(resp.Body)

	box := packr.NewBox("./assets")

	body, _ := box.MustBytes(url)
	avatar, _ := jpeg.Decode(bytes.NewReader(body))

	avatar = resize.Resize(16, 16, avatar, resize.NearestNeighbor)
	for x := x1; x < x1+16; x++ {
		for y := y1; y < y1+16; y++ {
			r, g, b, a := avatar.At(x-x1, y-y1).RGBA()

			/*fmt.Print("{");
			fmt.Print(uint8(r));
			fmt.Print(",");
			fmt.Print(uint8(g));
			fmt.Print(",");
			fmt.Print(uint8(b));
			fmt.Print(",");
			fmt.Print(uint8(a));
			fmt.Print("} ");*/

			if uint8(r) == 255 && uint8(g) == 255 && uint8(b) == 255 && uint8(a) == 255 {
				im.Set(x, y, color.RGBA{238, 238, 238, 255})
				continue
			}

			if theme == "dark" && uint8(r) == 0 && uint8(g) == 0 && uint8(b) == 0 && uint8(a) == 255 {
				im.Set(x, y, color.RGBA{34, 34, 34, 255})
				continue
			}

			//fmt.Print("  ")

			//if uint8(a) != 0 {
			//im.Set(x, y, color.RGBA{uint8(r), uint8(b), uint8(g), uint8(a)})
			im.Set(x, y, avatar.At(x-x1, y-y1))
			//}
			//fmt.Print(im.At(x,y));
		}
	}
	//fmt.Println();
}
func DownloadImages() {

	fmt.Println("Downloading Images ...")

	//preparing the background image
	clean = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{250, 90}})
	dark = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{250, 90}})

	border_color := color.RGBA{204, 204, 204, 255}

	//set the border
	for x := 0; x < 250; x++ {
		clean.Set(x, 0, border_color)
		clean.Set(x, 89, border_color)
		dark.Set(x, 0, border_color)
		dark.Set(x, 89, border_color)
	}
	for y := 0; y < 90; y++ {
		clean.Set(0, y, border_color)
		clean.Set(249, y, border_color)
		dark.Set(0, y, border_color)
		dark.Set(249, y, border_color)
	}

	//set background
	for x := 1; x < 249; x++ {
		for y := 1; y < 89; y++ {
			clean.Set(x, y, color.RGBA{238, 238, 238, 255})
			dark.Set(x, y, color.RGBA{34, 34, 34, 255})
		}
	}
	// var resp *http.Response

	/*img,_ := os.Open("images/github_dark.jpeg")
	avatar, _,_ := image.Decode(img)*/

	//set the github icon
	FillIcon(clean, 97, 7, "github.jpeg", "")
	FillIcon(dark, 97, 7, "github_dark.jpeg", "dark")

	//set repo icon
	FillIcon(clean, 97, 28, "repo.jpeg", "")
	FillIcon(dark, 97, 28, "repo_dark.jpeg", "dark")

	//set followers icon
	FillIcon(clean, 97, 48, "people.jpeg", "")
	FillIcon(dark, 97, 48, "people_dark.jpeg", "dark")

	//set fork icon
	FillIcon(clean, 97, 68, "fork.jpeg", "")
	FillIcon(dark, 97, 68, "fork_dark.jpeg", "dark")

	//set gist icon
	FillIcon(clean, 173, 48, "gist.jpeg", "")
	FillIcon(dark, 173, 48, "gist_dark.jpeg", "dark")

	//set star icon
	FillIcon(clean, 173, 28, "star.jpeg", "")
	FillIcon(dark, 173, 28, "star_dark.jpeg", "dark")

	// url := "http://cdn.steelhousemedia.com/files/docs/Creative/fonts/All%20Fonts/Arial%20Bold.ttf"
	// filename := "arialbd.ttf"
	// resp, _ = http.Get(url)

	// file, _ := os.Create(filename)
	// defer file.Close()

	// io.Copy(file, resp.Body)

	fmt.Println("Done")

}

func FetchCounts(username string) (string, string) {

	resp := HttpGet("https://api.github.com/users/" + username + "/repos")

	body, _ := ioutil.ReadAll(resp.Body)
	var data1 []map[string]interface{}
	json.Unmarshal(body, &data1)

	fork_count, star_count := 0, 0
	for v := range data1 {
		temp, _ := strconv.Atoi(fmt.Sprint(data1[v]["forks_count"]))
		fork_count = fork_count + temp

		temp, _ = strconv.Atoi(fmt.Sprint(data1[v]["stargazers_count"]))
		star_count = star_count + temp
	}
	return strconv.Itoa(fork_count), strconv.Itoa(star_count)
}

func HttpGet(url string) *http.Response {
	var resp *http.Response

	box := packr.NewBox("./secrets")

	body, err := box.MustString("access_token.txt")

	if err == nil {
		url = url + "?access_token=" + body[:len(body)-1]
	}

	resp, err = http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp
}
