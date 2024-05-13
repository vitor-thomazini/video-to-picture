package domain

import (
	"image"
	"image/draw"
)

type Quadrant struct {
	startPosition image.Point
	img           image.Image
	rectangle     image.Rectangle
}

func NewStartQuadrant(img image.Image, style Style) Quadrant {
	return Quadrant{
		img: img,
		startPosition: image.Point{
			X: img.Bounds().Min.X + style.MarginX,
			Y: img.Bounds().Min.Y + style.MarginY,
		},
	}
}

func NewStartQuadrantFromRect(rectangle image.Rectangle, img image.Image, style Style) Quadrant {
	return Quadrant{
		img: img,
		startPosition: image.Point{
			X: img.Bounds().Min.X + style.MarginX,
			Y: rectangle.Max.Y + style.PaddingY,
		},
	}
}

func NewMiddleQuadrant(rectangle image.Rectangle, img image.Image, style Style) Quadrant {
	return Quadrant{
		img: img,
		startPosition: image.Point{
			X: rectangle.Max.X + style.PaddingX,
			Y: rectangle.Min.Y,
		},
	}
}

func NewSecondQuadrant(img image.Image, style Style) Quadrant {
	return Quadrant{
		img: img,
		startPosition: image.Point{
			X: img.Bounds().Min.X + style.MarginX,
			Y: img.Bounds().Min.Y + style.MarginY,
		},
	}
}

func (q *Quadrant) Draw(background *image.RGBA) {
	q.rectangle = image.Rectangle{
		Min: q.startPosition,
		Max: q.startPosition.Add(q.img.Bounds().Size()),
	}
	draw.Src.Draw(background, q.rectangle, q.img, image.Point{Y: 0, X: 0})
}

func (q Quadrant) Rectangle() image.Rectangle {
	return q.rectangle
}
