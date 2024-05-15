package application

import (
	"fmt"
	"log"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

const (
	MAX_VALUE_FRAMES_COUNTER = 20
)

type FrameController struct {
	framesCounter         int
	framesByScreenCounter int
	latestText            string
}

func NewFrameController() FrameController {
	return FrameController{
		framesCounter:         0,
		framesByScreenCounter: 0,
	}
}

func (c *FrameController) Wait() bool {
	if c.framesCounter == MAX_VALUE_FRAMES_COUNTER {
		c.framesCounter = 0
		return false
	}
	c.framesCounter += 1
	return true
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
	SrcFilepath string
	DstFilepath string
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

func (w *GetFrameFromWhatsappVideo) Execute(params GetFrameFromWhatsappVideoParams) {
	video := w.captureVideo(params.SrcFilepath)
	imageToText := NewGetTextFromImage()
	defer imageToText.CloseTransaction()

	bkgList := make(map[string][]domain.Resource)

	c := 0

	for video.IsOpened() {
		img, notExistsImg := w.readVideo(&video)
		if notExistsImg {
			break
		}

		if w.frameController.Wait() {
			continue
		}

		// p := image.Point{X: domain.BACKGROUND_WIDTH, Y: domain.BACKGROUND_HEIGHT}
		// dst := gocv.NewMat()
		// gocv.Resize(*img, &dst, maxPoint, 0.1, 0.1, gocv.InterpolationLinear)
		// i, _ := dst.ToImage()

		frame := domain.NewFrame(*img).
			Resize(domain.BACKGROUND_HEIGHT, domain.BACKGROUND_WIDTH)

		texts := imageToText.Execute(frame)

		// fmt.Println("texts", texts)
		for _, text := range texts {
			w.frameController.InitLatestText(text)
			if w.frameController.latestText != text {
				// Cria novo bkg caso esteja cheio
				bkgList = domain.DrawAndUpdateResources(frame, bkgList, w.frameController.LatestText())
				w.frameController.UpdateLatestText(text)
				// Pode adicionar algo neste background
				bkgList, _ = domain.AppendLastResourceToResourceMap(bkgList, w.frameController.LatestText())
			} else {
				// Cria novo bkg caso esteja cheio
				bkgList, _ = domain.AppendLastResourceToResourceMap(bkgList, w.frameController.LatestText())
			}
		}

		fmt.Printf("text: %s, len: %d\n", w.frameController.LatestText(), len(bkgList[w.frameController.LatestText()]))

		bkgList = domain.DrawAndUpdateResources(frame, bkgList, w.frameController.LatestText())

		c += 1

		if c > 50 {
			break
		}
	}

	fmt.Println(bkgList)

	w.saveToPdf.Execute(params.DstFilepath, bkgList)
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
