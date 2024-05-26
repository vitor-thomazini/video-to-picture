import numpy
import cv2 as cv

from domain.position import Position

BACKGROUND_HEIGHT = 3508
BACKGROUND_WIDTH  = 2480

class Frame:
    def __init__(self, image: numpy.ndarray) -> None:
        self.image = image 

    def resize(self):
        p = self.__position()
        return cv.resize(self.image, p.asTuple())

    def __position(self) -> Position:
        x = self.image.shape[1]
        y = self.image.shape[0]
        dx = x - (BACKGROUND_WIDTH / 2)
        dy = x - (BACKGROUND_HEIGHT / 2)
        if (dx > dy) and (dx < 0 or dy < 0):
            x = x - dx
            y = y - (abs(dx) * (y / x))
            return Position(x, y)
        elif (dx < dy) and (dx < 0 or dy < 0):
            x = x - (abs(dy) * (x / y))
            y = y - dy
            return Position(x, y)
        else:
            return Position(x, y)




    # imgX := float64(f.img.Size()[1])
	# imgY := float64(f.img.Size()[0])
	# dy := imgY - (float64(panelHeight) / 2)
	# dx := imgX - (float64(panelWidth) / 2)

	# if (dx > dy) && (dx < 0 || dy < 0) {
	# 	return image.Point{
	# 		X: int(imgX-dx) - f.style.MarginX,
	# 		Y: int(imgY-(math.Abs(dx)*imgY/imgX)) - f.style.MarginY,
	# 	}
	# } else if (dx < dy) && (dx < 0 || dy < 0) {
	# 	return image.Point{
	# 		X: int(imgX-(dy*imgX/imgY)) - f.style.MarginX,
	# 		Y: int(imgY-dy) - f.style.MarginY,
	# 	}
	# } else {
	# 	// TODO: implement resizing both dimension, based on A4 halt
	# 	return image.Point{}
	# }