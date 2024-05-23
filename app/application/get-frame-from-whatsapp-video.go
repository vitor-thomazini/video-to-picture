package application

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

type GetFrameFromWhatsappVideoParams struct {
}

type GetFrameFromWhatsappVideo struct {
	frameController domain.FrameController
	saveToPdf       *SaveToPdf
}

func NewGetFrameFromWhatsappVideo() GetFrameFromWhatsappVideo {
	return GetFrameFromWhatsappVideo{

		saveToPdf: NewSaveToPdf(),
	}
}

func (w *GetFrameFromWhatsappVideo) Execute(srcFile string, dstDir string) {
	video := w.captureVideo(srcFile)
	imageToText := NewGetTextFromImage()
	defer imageToText.CloseTransaction()

	saveStorage := NewSaveStorageBatch()
	storage := NewLoadStorageBatch().Execute()
	w.frameController = domain.NewFrameController(storage.LatestIndex, 1000)

	bkgList := storage.GetData()
	for video.IsOpened() {
		img, notExistsImg := w.readVideo(&video)
		if notExistsImg {
			break
		}

		if w.frameController.Wait() {
			continue
		}

		texts := imageToText.Execute(*img)
		frame, _ := domain.NewFrame(*img).
			Resize(domain.BACKGROUND_HEIGHT, domain.BACKGROUND_WIDTH)

		fmt.Printf("%d - %s\n", w.frameController.Counter(), strings.Join(texts, " "))
		for _, text := range texts {
			bkgList = domain.DrawAndUpdateResources(frame, bkgList, text)
		}

		if len(texts) > 0 {
			sort.Slice(texts, func(i, j int) bool {
				return texts[j] < texts[i]
			})
			w.frameController.UpdateLatestText(texts[0])
		} else {
			bkgList = domain.DrawAndUpdateResources(frame, bkgList, w.frameController.LatestText())
		}

		if w.frameController.Counter() >= w.frameController.EndInIndexFrame() {
			saveStorage.Execute(w.frameController.LatestText(), w.frameController.Counter(), bkgList)
			break
		}

	}

	fmt.Println(bkgList)

	// w.saveToPdf.Execute(dstDir, bkgList)
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
