package main

import (
	"github.com/vitor-thomazini/video-to-picture/app/application"
)

func main() {
	a := application.NewGetFrameFromWhatsappVideo()
	a.Execute(application.GetFrameFromWhatsappVideoParams{
		SrcFilepath: "video.mp4",
		DstFilepath: "result",
	})

	// cmd.Execute()
}
