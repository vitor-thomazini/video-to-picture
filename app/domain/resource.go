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
	image          *image.RGBA
	pictureCounter int
	background     image.Rectangle
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
		image:          img,
		pictureCounter: 0,
	}
}

func (r *Resource) UpdateBackground(background image.Rectangle) {
	r.background = background
}

func (r Resource) Image() *image.RGBA {
	return r.image
}

func (r Resource) Background() *image.Rectangle {
	return &r.background
}

func (q Resource) IsFirstPictureToFirstRow() bool {
	return q.pictureCounter == 0
}

func (q Resource) IsFirstPictureToAnyRow() bool {
	return (q.pictureCounter != 0) && (q.pictureCounter%SPLIT_FRAMES_SCREEN == 0)
}

func (q Resource) IsLastPicture() bool {
	return q.pictureCounter >= TOTAL_FRAMES_BY_SCREEN
}

func (q *Resource) IncPictureCounter() {
	q.pictureCounter += 1
}
