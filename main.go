package main

import (
	"github.com/vitor-thomazini/video-to-picture/app/application"
)

func main() {
	a := application.NewConvertVideoToImage("video.mp4", "result/images.pdf", 10)
	a.Convert()

	// cmd.Execute()
}
