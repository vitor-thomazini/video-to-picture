package application

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/vitor-thomazini/video-to-picture/app/domain"
	"gocv.io/x/gocv"
)

type Counter struct {
	frame    int
	quadrant int
	page     int
}

func NewCounter() Counter {
	return Counter{
		frame:    0,
		quadrant: 0,
		page:     -1,
	}
}

func (c *Counter) InitQuadrant() {
	c.quadrant = 0
}

func (c *Counter) IncQuadrant() {
	c.quadrant += 1
}

func (c *Counter) IncFrame() {
	c.frame += 1
}

func (c *Counter) IncPage() {
	c.page += 1
}

type ConvertVideoToImage struct {
	filepath         string
	totalFrameInPage int
	panel            image.Rectangle
	counter          Counter
	background       domain.Background
	backgrounds      []domain.Background
	style            domain.Style
}

func NewConvertVideoToImage(filepath string) ConvertVideoToImage {
	return ConvertVideoToImage{
		filepath:         filepath,
		totalFrameInPage: 6,
		style: domain.Style{
			MarginX:  48,
			MarginY:  28,
			PaddingX: 64,
			PaddingY: 64,
		},
		counter:     NewCounter(),
		backgrounds: make([]domain.Background, 0),
	}
}

func (c *ConvertVideoToImage) Convert() {
	video, err := gocv.VideoCaptureFile(c.filepath)
	if err != nil {
		log.Fatalln("File not found")
	}

	img := gocv.NewMat()
	for video.IsOpened() {
		canRead := video.Read(&img)
		if !canRead {
			break
		}

		background := c.initializePage()
		frame := domain.NewFrame(img, background).Image()
		c.drawQuadrants(frame)

		c.counter.IncQuadrant()
		c.counter.IncFrame()
		if c.counter.frame > 6 {
			break
		}

	}

	c.saveImages()
}

func (c *ConvertVideoToImage) initializePage() domain.Background {
	if c.isStartedPage() {
		c.background = domain.NewA4Dimension()
		c.backgrounds = append(c.backgrounds, c.background)
		c.counter.InitQuadrant()
		c.counter.IncPage()
		fmt.Println(c.counter, "initializePage")
	}
	return c.background
}

func (c *ConvertVideoToImage) drawQuadrants(frame image.Image) {
	var quadrant domain.Quadrant

	if c.isFirsQuadrantForFirstRow() {
		quadrant = domain.NewStartQuadrant(frame, c.style)
		quadrant.Draw(c.backgrounds[c.counter.page].Image())
		c.panel = quadrant.Rectangle()
	} else if c.isFirsQuadrantForAnyRow() {
		quadrant = domain.NewStartQuadrantFromRect(c.panel, frame, c.style)
		quadrant.Draw(c.backgrounds[c.counter.page].Image())
		c.panel = quadrant.Rectangle()
	} else {
		quadrant = domain.NewMiddleQuadrant(c.panel, frame, c.style)
		quadrant.Draw(c.backgrounds[c.counter.page].Image())
		c.panel = quadrant.Rectangle()
	}

}

func (c ConvertVideoToImage) saveImages() {
	for idx, background := range c.backgrounds {
		w, _ := os.Create(fmt.Sprintf("result/%d.jpg", idx))
		defer w.Close()
		jpeg.Encode(w, background.Image(), &jpeg.Options{Quality: 90})
	}
}

func (c ConvertVideoToImage) isStartedPage() bool {
	return (c.counter.frame == 0) || (c.counter.frame%c.totalFrameInPage == 0)
}

func (c ConvertVideoToImage) isFirsQuadrantForFirstRow() bool {
	return c.counter.quadrant == 0
}

func (c ConvertVideoToImage) isFirsQuadrantForAnyRow() bool {
	return (c.counter.quadrant != 0) && (c.counter.quadrant%3 == 0)
}
