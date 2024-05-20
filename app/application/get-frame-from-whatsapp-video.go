package application

import (
	"fmt"
	"log"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

const (
	MAX_VALUE_FRAMES_COUNTER = 12
)

type FrameController struct {
	framesCounter int
	latestText    string
}

func NewFrameController() FrameController {
	return FrameController{
		framesCounter: 0,
	}
}

func (c *FrameController) Wait() bool {
	wait := c.framesCounter%MAX_VALUE_FRAMES_COUNTER != 0
	c.framesCounter += 1
	return wait
}

func (c *FrameController) InitLatestText(text string) {
	if c.latestText == "" {
		c.latestText = text
	}
}

func (c *FrameController) UpdateLatestText(text string) {
	if c.latestText != "none" {
		c.latestText = text
	}
}

func (c FrameController) LatestText() string {
	return c.latestText
}

type GetFrameFromWhatsappVideoParams struct {
}

type GetFrameFromWhatsappVideo struct {
	frameController FrameController
	saveToPdf       *SaveToPdf
}

func NewGetFrameFromWhatsappVideo() GetFrameFromWhatsappVideo {
	return GetFrameFromWhatsappVideo{
		frameController: NewFrameController(),
		saveToPdf:       NewSaveToPdf(),
	}
}

func (w *GetFrameFromWhatsappVideo) Execute(srcFile string, dstDir string) {
	video := w.captureVideo(srcFile)
	imageToText := NewGetTextFromImage()
	defer imageToText.CloseTransaction()

	bkgList := make(map[string][]domain.Resource)
	for video.IsOpened() {
		img, notExistsImg := w.readVideo(&video)
		if notExistsImg {
			break
		}

		if w.frameController.Wait() {
			continue
		}

		texts := imageToText.Execute(*img)
		frame := domain.NewFrame(*img).
			Resize(domain.BACKGROUND_HEIGHT, domain.BACKGROUND_WIDTH)

		for _, text := range texts {
			w.frameController.InitLatestText(text)
			if w.frameController.latestText != text && len(texts) > 1 {
				bkgList = domain.DrawAndUpdateResources(frame, bkgList, w.frameController.LatestText())
			}

			if w.frameController.latestText != text {
				w.frameController.UpdateLatestText(text)
			}

			bkgList, _ = domain.AppendLastResourceToResourceMap(bkgList, w.frameController.LatestText())
		}

		bkgList = domain.DrawAndUpdateResources(frame, bkgList, w.frameController.LatestText())
	}

	fmt.Println(bkgList)

	w.saveToPdf.Execute(dstDir, bkgList)
}

func (w GetFrameFromWhatsappVideo) captureVideo(srcFilepath string) gocv.VideoCapture {
	video, err := gocv.VideoCaptureFile(srcFilepath)
	if err != nil {
		log.Fatalln("File not found")
	}

	return *video
}

func (w GetFrameFromWhatsappVideo) readVideo(video *gocv.VideoCapture) (*gocv.Mat, bool) {
	img := gocv.NewMat()
	exists := video.Read(&img)
	if !exists {
		return nil, !exists
	}
	return &img, !exists
}
