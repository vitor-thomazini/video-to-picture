package domain

import "image"

const (
	TOTAL_FRAMES_BY_SCREEN = 6
	SPLIT_FRAMES_SCREEN    = 3
)

type Resource struct {
	image          *image.RGBA
	pictureCounter int
	background     image.Rectangle
}

func NewResource(image *image.RGBA) Resource {
	return Resource{
		image:          image,
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
