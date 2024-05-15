package domain

import (
	"image"
	"log"
	"math"

	"gocv.io/x/gocv"
)

type Frame struct {
	img   gocv.Mat
	style Style
}

func NewFrame(img gocv.Mat) *Frame {
	return &Frame{
		img: img,
		style: Style{
			MarginX: 60,
			MarginY: 60,
		},
	}
}

func (f *Frame) Resize(panelWidth int, panelHeight int) image.Image {
	dst := gocv.NewMat()
	maxPoint := f.getPointFromPanel(panelWidth, panelHeight)
	gocv.Resize(f.img, &dst, maxPoint, 0.1, 0.1, gocv.InterpolationLinear)

	img, err := dst.ToImage()
	if err != nil {
		log.Fatalln(err)
	}
	return img
}

func (f *Frame) getPointFromPanel(panelWidth int, panelHeight int) image.Point {
	imgX := float64(f.img.Size()[1])
	imgY := float64(f.img.Size()[0])
	dy := imgY - (float64(panelHeight) / 2)
	dx := imgX - (float64(panelWidth) / 2)

	if (dx > dy) && (dx < 0 || dy < 0) {
		return image.Point{
			X: int(imgX-dx) - f.style.MarginX,
			Y: int(imgY-(math.Abs(dx)*imgY/imgX)) - f.style.MarginY,
		}
	} else if (dx < dy) && (dx < 0 || dy < 0) {
		return image.Point{
			X: int(imgX-(dy*imgX/imgY)) - f.style.MarginX,
			Y: int(imgY-dy) - f.style.MarginY,
		}
	} else {
		// TODO: implement resizing both dimension, based on A4 halt
		return image.Point{}
	}
}
