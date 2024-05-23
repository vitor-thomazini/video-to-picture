package domain

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	TOTAL_FRAMES_BY_SCREEN = 6
	SPLIT_FRAMES_SCREEN    = 3
)

const (
	BACKGROUND_HEIGHT = 3508
	BACKGROUND_WIDTH  = 2480
)

type Resource struct {
	Image          *image.RGBA
	PictureCounter int
	Background     image.Rectangle
}

func NewResource() Resource {
	minPoint := image.Point{X: 0, Y: 0}
	maxPoint := image.Point{X: BACKGROUND_WIDTH, Y: BACKGROUND_HEIGHT}

	img := image.NewRGBA(image.Rectangle{
		Min: minPoint,
		Max: maxPoint,
	})

	whiteColor := image.Uniform{C: color.RGBA{255, 255, 255, 1.0}}
	draw.Src.Draw(img, img.Bounds(), &whiteColor, image.Point{})

	return Resource{
		Image:          img,
		PictureCounter: 0,
	}
}

func (r *Resource) UpdateBackground(background image.Rectangle) {
	r.Background = background
}

func (q Resource) IsFirstPictureToFirstRow() bool {
	return q.PictureCounter == 0
}

func (q Resource) IsFirstPictureToAnyRow() bool {
	return (q.PictureCounter != 0) && (q.PictureCounter%SPLIT_FRAMES_SCREEN == 0)
}

func (q Resource) IsLastPicture() bool {
	return q.PictureCounter >= TOTAL_FRAMES_BY_SCREEN
}

func (q *Resource) IncPictureCounter() {
	q.PictureCounter += 1
}
