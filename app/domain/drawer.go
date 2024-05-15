package domain

import (
	"image"
	"image/draw"
)

type Drawer struct {
	startPosition image.Point
	img           image.Image
}

func NewDrawerStartPoint(img image.Image, style Style) Drawer {
	return Drawer{
		img: img,
		startPosition: image.Point{
			X: img.Bounds().Min.X + style.MarginX,
			Y: img.Bounds().Min.Y + style.MarginY,
		},
	}
}

func NewDrawerStartPointFromRect(rectangle image.Rectangle, img image.Image, style Style) Drawer {
	return Drawer{
		img: img,
		startPosition: image.Point{
			X: img.Bounds().Min.X + style.MarginX,
			Y: rectangle.Max.Y + style.PaddingY,
		},
	}
}

func NewDrawerMiddlePoint(rectangle image.Rectangle, img image.Image, style Style) Drawer {
	return Drawer{
		img: img,
		startPosition: image.Point{
			X: rectangle.Max.X + style.PaddingX,
			Y: rectangle.Min.Y,
		},
	}
}

func (q Drawer) CalculatePanel() image.Rectangle {
	return image.Rectangle{
		Min: q.startPosition,
		Max: q.startPosition.Add(q.img.Bounds().Size()),
	}
}

func (q Drawer) DrawIn(bkg *image.RGBA, panel *image.Rectangle) {
	draw.Src.Draw(bkg, *panel, q.img, image.Point{Y: 0, X: 0})
}
