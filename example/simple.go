package main

import (
	"fmt"
	"log"

	youtube "github.com/jameshwc/Youtube-Subtitle-Downloader"
)

func main() {
	url := "http://youtu.be/5MgBikgcWnY"
	subtitle, err := youtube.Download(url)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range subtitle.Lines {
		fmt.Println(line.Start, line.End, line.Text)
	}
}
