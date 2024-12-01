package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type BingImage struct {
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Url           string `json:"url"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyright_link"`
}

var Market map[string]string

func main() {
	market := getMarket()
	imageData := todayBingImageData(market)
	imagePath := downloadBingImage(imageData.Url)

	log.Println("Setting wallpaper from bing [" + market + "]")
	err := exec.Command("osascript", "-e", "tell application \"System Events\" to tell every desktop to set picture to POSIX file \""+imagePath+"\"").Run()
	check(err)
}

func getMarket() string {
	Market = map[string]string{
		"us":     "en-US",
		"es":     "es-ES",
		"jp":     "ja-JP",
		"random": "random",
	}

	mktFlag := flag.String("mkt", "random", "Market to get the image from")
	flag.Parse()

	mkt := strings.ToLower(*mktFlag)
	if _, ok := Market[mkt]; !ok {
		mkt = "random"
	}
	return Market[mkt]
}

func todayBingImageData(market string) BingImage {
	resp, err := http.Get("https://bing.biturl.top/?resolution=UHD&format=json&index=0&mkt=" + market)
	check(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	check(err)

	var data BingImage
	json.Unmarshal(body, &data)
	return data
}

func downloadBingImage(bingImageUrl string) string {
	resp, err := http.Get(bingImageUrl)
	check(err)

	defer resp.Body.Close()
	rawImage, err := io.ReadAll(resp.Body)
	check(err)

	imageId := resp.Request.URL.Query().Get("id")
	imageFile, err := os.CreateTemp("", imageId)
	check(err)

	imageFile.Write(rawImage)
	return imageFile.Name()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
