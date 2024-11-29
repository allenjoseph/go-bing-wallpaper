package main

import (
	"encoding/json"
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
	market := getMarket(readStdin())
	imageData := todayBingImageData(market)
	imagePath := downloadBingImage(imageData.Url)

	log.Println("Setting wallpaper from bing [" + market + "]")
	err := exec.Command("osascript", "-e", "tell application \"System Events\" to tell every desktop to set picture to POSIX file \""+imagePath+"\"").Run()
	check(err)
}

func readStdin() string {
	stdin, err := io.ReadAll(os.Stdin)
	check(err)

	return strings.TrimSpace(string(stdin))
}

func getMarket(region string) string {
	Market = map[string]string{
		"cn":     "zh-CN",
		"us":     "en-US",
		"jp":     "ja-JP",
		"au":     "en-AU",
		"uk":     "en-GB",
		"ge":     "de-DE",
		"nz":     "en-NZ",
		"ca":     "en-CA",
		"random": "random",
	}

	market := Market["random"]
	if _, ok := Market[region]; ok {
		market = Market[region]
	}
	return market
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
