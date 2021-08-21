package main

import (
	"context"
	"log"
	"time"

	ytdl "github.com/kkdai/youtube/v2/downloader"
)

const tempFile = "video.mp4"
const videoID = "NQ9RtLrapzc"
const timeout = 3

func main() {
	dl := &ytdl.Downloader{}
	video, err := dl.GetVideo(videoID)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		ctx, _ := context.WithTimeout(context.Background(), timeout*time.Second)
		err := dl.Download(ctx, video, &video.Formats[0], tempFile)
		if err == nil {
			log.Printf("scientific browsing still good, true internet")
		} else {
			log.Printf("scientific browsing not good, download failed %s, retrying", err.Error())
		}

		<-time.After(timeout * time.Second)
	}
}
