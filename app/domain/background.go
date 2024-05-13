package domain

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"gocv.io/x/gocv"
)

type Background struct {
	height int
	width  int
	img    *image.RGBA
}

func NewA4Dimension() Background {
	background := Background{
		height: 3508,
		width:  2480,
	}
	background.img = image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: background.width,
			Y: background.height,
		},
	})

	whiteColor := color.RGBA{255, 255, 255, 1.0}
	draw.Draw(
		background.img,
		background.img.Bounds(),
		&image.Uniform{whiteColor},
		image.Point{},
		draw.Src,
	)

	return background
}

func (d Background) Width() int {
	return d.width
}

func (d Background) Height() int {
	return d.height
}

func (d Background) GetPointResized(img gocv.Mat, style Style) image.Point {
	imgX := float64(img.Size()[1])
	imgY := float64(img.Size()[0])
	dy := imgY - (float64(d.height) / 2)
	dx := imgX - (float64(d.width) / 2)

	if (dx > dy) && (dx < 0 || dy < 0) {
		return image.Point{
			X: int(imgX-dx) - style.MarginX,
			Y: int(imgY-(math.Abs(dx)*imgY/imgX)) - style.MarginY,
		}
	} else if (dx < dy) && (dx < 0 || dy < 0) {
		return image.Point{
			X: int(imgX-(dy*imgX/imgY)) - style.MarginX,
			Y: int(imgY-dy) - style.MarginY,
		}
	} else {
		// TODO: implement resizing both dimension, based on A4 halt
		return image.Point{}
	}
}

func (d Background) Image() *image.RGBA {
	return d.img
}

func (d Background) DrawImage(img image.Image) {}
