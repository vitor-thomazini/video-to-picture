package application

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

const (
	MAX_VALUE_FRAMES_COUNTER = 20
	BACKGROUND_HEIGHT        = 3508
	BACKGROUND_WIDTH         = 2480
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

func (c FrameController) IsFirstFrameInScreen() bool {
	return c.framesByScreenCounter == 0
}

func (c FrameController) CanCreateBackground(text string) bool {
	return (c.IsFirstFrameInScreen() && c.latestText == text) || (c.latestText != text)
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

		frame := domain.NewFrame(*img).
			Resize(BACKGROUND_HEIGHT, BACKGROUND_WIDTH)

		texts := imageToText.Execute(frame)
		// fmt.Println("texts", texts)
		for _, text := range texts {
			w.frameController.InitLatestText(text)
			if w.frameController.latestText != text {
				bkg := w.lastBkg(bkgList[w.frameController.LatestText()])
				// Cria novo bkg caso esteja cheio
				if bkg == (domain.Resource{}) || bkg.IsLastPicture() {
					bkg = w.createBackground()
					bkgList[w.frameController.LatestText()] = append(bkgList[w.frameController.LatestText()], bkg)
				}
				bkg = w.drawIn(frame, bkg)
				bkgList[w.frameController.LatestText()] = w.updateLatestBkg(bkgList[w.frameController.LatestText()], bkg)

				w.frameController.UpdateLatestText(text)

				bkg = w.lastBkg(bkgList[w.frameController.LatestText()])
				// Pode adicionar algo neste background
				if bkg == (domain.Resource{}) || bkg.IsLastPicture() {
					bkg = w.createBackground()
					bkgList[w.frameController.LatestText()] = append(bkgList[w.frameController.LatestText()], bkg)
				}
			} else {
				bkg := w.lastBkg(bkgList[w.frameController.LatestText()])
				// Cria novo bkg caso esteja cheio
				if bkg == (domain.Resource{}) || bkg.IsLastPicture() {
					bkg = w.createBackground()
					bkgList[w.frameController.LatestText()] = append(bkgList[w.frameController.LatestText()], bkg)
				}
			}
		}

		fmt.Printf("text: %s, len: %d\n", w.frameController.LatestText(), len(bkgList[w.frameController.LatestText()]))

		bkg := w.lastBkg(bkgList[w.frameController.LatestText()])
		if bkg == (domain.Resource{}) || bkg.IsLastPicture() {
			bkg = w.createBackground()
			bkgList[w.frameController.LatestText()] = append(bkgList[w.frameController.LatestText()], bkg)
		}
		bkg = w.drawIn(frame, bkg)
		bkgList[w.frameController.LatestText()] = w.updateLatestBkg(bkgList[w.frameController.LatestText()], bkg)

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

func (w GetFrameFromWhatsappVideo) createBackground() domain.Resource {
	minPoint := image.Point{X: 0, Y: 0}
	maxPoint := image.Point{X: BACKGROUND_WIDTH, Y: BACKGROUND_HEIGHT}

	img := image.NewRGBA(image.Rectangle{
		Min: minPoint,
		Max: maxPoint,
	})

	whiteColor := color.RGBA{255, 255, 255, 1.0}
	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{whiteColor},
		image.Point{},
		draw.Src,
	)
	return domain.NewResource(img)
}

func (w *GetFrameFromWhatsappVideo) drawIn(frame image.Image, bkg domain.Resource) domain.Resource {
	var drawer domain.Drawer
	style := domain.Style{
		MarginX:  112,
		MarginY:  28,
		PaddingX: 8,
		PaddingY: 8,
	}
	if bkg.IsFirstPictureToFirstRow() {
		drawer = domain.NewDrawerStartPoint(frame, style)
	} else if bkg.IsFirstPictureToAnyRow() {
		drawer = domain.NewDrawerStartPointFromRect(*bkg.Background(), frame, style)
	} else {
		drawer = domain.NewDrawerMiddlePoint(*bkg.Background(), frame, style)
	}

	bkg.UpdateBackground(drawer.CalculatePanel())
	drawer.DrawIn(bkg.Image(), bkg.Background())
	bkg.IncPictureCounter()
	return bkg
}

func (w GetFrameFromWhatsappVideo) lastBkg(bkgList []domain.Resource) domain.Resource {
	if len(bkgList) == 0 {
		return domain.Resource{}
	}
	lastIndex := len(bkgList) - 1
	return bkgList[lastIndex]
}

func (w GetFrameFromWhatsappVideo) updateLatestBkg(bkgList []domain.Resource, bkg domain.Resource) []domain.Resource {
	if len(bkgList) == 0 {
		return []domain.Resource{}
	}
	lastIndex := len(bkgList) - 1
	bkgList[lastIndex] = bkg
	return bkgList
}
