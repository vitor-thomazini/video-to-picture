package application

import (
	"fmt"
	"image"
	"log"

	"github.com/signintech/gopdf"
	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

type ConvertVideoToImage struct {
	videoFilepath    string
	pdfFilepath      string
	fps              int
	totalFrameInPage int
	panel            image.Rectangle
	counter          domain.Counter
	background       domain.Background
	backgrounds      []domain.Background
	style            domain.Style
}

func NewConvertVideoToImage(
	videoFilepath string,
	pdfFilepath string,
	fps int,
) ConvertVideoToImage {
	return ConvertVideoToImage{
		videoFilepath:    videoFilepath,
		pdfFilepath:      pdfFilepath,
		fps:              fps,
		totalFrameInPage: 6,
		style: domain.Style{
			MarginX:  112,
			MarginY:  28,
			PaddingX: 8,
			PaddingY: 8,
		},
		counter:     domain.NewCounter(),
		backgrounds: make([]domain.Background, 0),
	}
}

func (c *ConvertVideoToImage) Convert() {
	video, err := gocv.VideoCaptureFile(c.videoFilepath)
	if err != nil {
		log.Fatalln("File not found")
	}

	img := gocv.NewMat()
	for video.IsOpened() {
		canRead := video.Read(&img)
		if !canRead {
			break
		}
		if c.counter.CanIgnore(c.fps) {
			continue
		}

		background := c.initializePage()
		frame := domain.NewFrame(img, background).Image()
		c.drawQuadrants(frame)

		c.counter.IncQuadrant()
		c.counter.IncFrame()
	}

	c.save()
}

func (c *ConvertVideoToImage) initializePage() domain.Background {
	if c.counter.IsStartPage(c.totalFrameInPage) {
		c.background = domain.NewA4Dimension()
		c.backgrounds = append(c.backgrounds, c.background)
		c.counter.InitQuadrant()
		c.counter.IncPage()
	}
	return c.background
}

func (c *ConvertVideoToImage) drawQuadrants(frame image.Image) {
	var quadrant domain.Quadrant

	if c.counter.IsFirstQuadrantForFirstRow() {
		quadrant = domain.NewStartQuadrant(frame, c.style)
	} else if c.counter.IsFirstQuadrantForAnyRow() {
		quadrant = domain.NewStartQuadrantFromRect(c.panel, frame, c.style)
	} else {
		quadrant = domain.NewMiddleQuadrant(c.panel, frame, c.style)
	}

	quadrant.Draw(c.backgrounds[c.counter.Page()].Image())
	c.panel = quadrant.Rectangle()
}

func (c ConvertVideoToImage) save() {
	var pdf gopdf.GoPdf

	rect := &gopdf.Rect{
		W: 595, H: 842,
	}
	rect.PointsToUnits(gopdf.UnitPX)

	pdf.Start(gopdf.Config{
		PageSize: *rect,
		Unit:     gopdf.UnitPX,
	})

	fmt.Println(len(c.backgrounds))
	for _, background := range c.backgrounds {
		pdf.AddPage()
		pdf.ImageFrom(background.Image(), 0, 0, rect)
	}

	pdf.WritePdf(c.pdfFilepath)
}
