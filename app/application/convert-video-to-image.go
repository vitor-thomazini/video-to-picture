package application

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"gocv.io/x/gocv"
)

const (
	TOTAL_QUADRANTS = 7
	MARGIN_Y        = 28
	MARGIN_X        = 56
	PADDING_X       = 64
	PADDING_Y       = 64
)

type ImageSize struct {
	height int
	width  int
}

func NewA4ImageSize() ImageSize {
	return ImageSize{
		height: 3508,
		width:  2480,
	}
}

func (i ImageSize) GetFourthPartSize(img gocv.Mat) image.Point {
	fmt.Println(img.Size()[0], i.height/2) // ABS(2340-1754)=600 ---> Y = 2340-600=1754-marginY
	fmt.Println(img.Size()[1], i.width/2)  // ABS(1080-1240)=160  ---> X = 1080-(600x1080/2340)=803-marginX
	// if img.Size()[0] < (i.height/2) && img.Size()[1] < (i.width/2) {
	// 	return image.Point{
	// 		Y: img.Size()[0],
	// 		X: img.Size()[1],
	// 	}
	// }

	// return image.Point{
	// 	Y: int(((i.height / 2) * img.Size()[0]) / i.height),
	// 	X: int(((i.width / 2) * img.Size()[1]) / i.width),
	// }

	return image.Point{
		Y: 1754 - 60,
		X: 803 - 60,
	}

}

func NewFourthPartA4ImageSize() ImageSize {
	return ImageSize{
		height: 3508 / 2,
		width:  2480 / 2,
	}
}

type ConvertVideoToImage struct {
	filepath string
}

func NewConvertVideoToImage(filepath string) ConvertVideoToImage {
	return ConvertVideoToImage{
		filepath: filepath,
	}
}

func (c ConvertVideoToImage) Convert() {
	video, err := gocv.VideoCaptureFile(c.filepath)
	if err != nil {
		log.Fatalln("File not found")
	}

	img := gocv.NewMat()
	size := NewA4ImageSize()

	frame := 0
	currentQuadrant := 0

	// startPosition := image.Rectangle{}
	backgrounds := make([]*image.RGBA, 0)
	page := -1

	var rectangle image.Rectangle
	for video.IsOpened() {
		canRead := video.Read(&img)
		if canRead {
			dst := gocv.NewMat()
			resizePoint := size.GetFourthPartSize(img)
			gocv.Resize(img, &dst, resizePoint, 0.1, 0.1, gocv.InterpolationLinear)

			im, err := dst.ToImage()
			if err != nil {
				log.Fatalln(err)
			}

			if frame == 0 || frame%TOTAL_QUADRANTS == 0 {
				background := image.NewRGBA(image.Rectangle{
					Min: image.Point{X: 0, Y: 0},
					Max: image.Point{X: size.width, Y: size.height},
				})
				draw.Draw(background, background.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 1.0}}, image.Point{}, draw.Src)
				backgrounds = append(backgrounds, background)
				currentQuadrant = 0
				page += 1
				fmt.Println("Initialize Page", frame)
			}

			// ROW 1
			if currentQuadrant == 0 {
				startPosition := image.Point{
					X: im.Bounds().Min.X + MARGIN_X,
					Y: im.Bounds().Min.Y + MARGIN_Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
			} else if currentQuadrant == 1 {
				startPosition := image.Point{
					X: rectangle.Max.X + PADDING_X,
					Y: rectangle.Min.Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
			} else if currentQuadrant == 2 {
				startPosition := image.Point{
					X: rectangle.Max.X + PADDING_X,
					Y: rectangle.Min.Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
				// ROW 2
			} else if currentQuadrant == 3 {
				startPosition := image.Point{
					X: im.Bounds().Min.X + MARGIN_X,
					Y: rectangle.Max.Y + PADDING_Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
			} else if currentQuadrant == 4 {
				startPosition := image.Point{
					X: rectangle.Max.X + PADDING_X,
					Y: rectangle.Min.Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
			} else if currentQuadrant == 5 {
				startPosition := image.Point{
					X: rectangle.Max.X + PADDING_X,
					Y: rectangle.Min.Y,
				}
				rectangle = image.Rectangle{startPosition, startPosition.Add(im.Bounds().Size())}
				draw.Src.Draw(backgrounds[page], rectangle, im, image.Point{Y: 0, X: 0})
			}

			gocv.WaitKey(20)

			currentQuadrant += 1
			frame += 1

			if frame > 6 {
				break
			}
		} else {
			break
		}
	}

	fmt.Println("Total Pages", len(backgrounds))
	for _, background := range backgrounds {
		w, _ := os.Create(fmt.Sprintf("result/%d.jpg", page))
		defer w.Close()
		jpeg.Encode(w, background, &jpeg.Options{Quality: 90})
		fmt.Println("=====")
		page += 1
	}
}
