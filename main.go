package main

import (
	"flag"
	"log"

	"github.com/vitor-thomazini/video-to-picture/app/application"
)

func main() {
	videoFilepath := flag.String("video", "/Users/vitor/Library/CloudStorage/OneDrive-Personal/video1.mp4", "Video filepath")
	dstDir := flag.String("dir", "/Users/vitor/Library/CloudStorage/OneDrive-Personal/test", "Target directory")
	flag.Parse()

	if videoFilepath == nil || *videoFilepath == "" {
		log.Fatalln("Video not found")
	}

	if dstDir == nil || *dstDir == "" {
		log.Fatalln("Directory not found")
	}

	a := application.NewGetFrameFromWhatsappVideo()
	a.Execute(*videoFilepath, *dstDir)
}
