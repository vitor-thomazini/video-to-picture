package domain

import (
	"image"
	"log"

	"gocv.io/x/gocv"
)

type Frame struct {
	imgMat     gocv.Mat
	background Background
	img        image.Image
	style      Style
}

func NewFrame(imgMat gocv.Mat, background Background) Frame {
	frame := Frame{
		imgMat:     imgMat,
		background: background,
		style: Style{
			MarginX: 60,
			MarginY: 60,
		},
	}
	frame.resize()
	return frame
}

func (f Frame) Image() image.Image {
	return f.img
}

func (f *Frame) resize() {
	dst := gocv.NewMat()
	maxPoint := f.background.GetPointResized(f.imgMat, f.style)
	gocv.Resize(f.imgMat, &dst, maxPoint, 0.1, 0.1, gocv.InterpolationLinear)

	var err error
	f.img, err = dst.ToImage()
	if err != nil {
		log.Fatalln(err)
	}
}
