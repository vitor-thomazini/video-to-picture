import numpy as np

from domain.position import Position

BACKGROUND_HEIGHT = 3508
BACKGROUND_WIDTH  = 2480

TOTAL_FRAMES_BY_SCREEN = 6
SPLIT_FRAMES_SCREEN    = 3

class Resource:
    def __init__(self) -> None:
        self.images = []
        self.counter = 0

    def isLast(self) -> bool:
        return self.counter >= TOTAL_FRAMES_BY_SCREEN
    
    def add(self, image: np.ndarray) -> None:
        self.images = [*self.images, image]
         
    def concat(self, index:int, image: np.ndarray, axis=1) -> None:
        self.images[index] = np.concatenate((self.images[index], image), axis)

    def fullImage(self) -> np.ndarray:
        # result = np.full((4680, 3240), 255)
        result = np.full((4680, 3240, 3), 255)
        for i in range(len(self.images)):
            y, x, _ = self.images[i].shape
            if i == 0: 
                dx = abs(x - 3240)
                if dx != 0:
                    self.images[i] = np.concatenate((self.images[i], np.full((2340, dx, 3), 255)), axis=1)
                y, x, _ = self.images[i].shape
                result[0:y][0:x] = self.images[i]
            else:
                y2, x2, _ = self.images[i-1].shape
                dx = abs(x - 3240)
                if dx != 0:
                    self.images[i] = np.concatenate((self.images[i], np.full((2340, dx, 3), 255)), axis=1)
                y, x, _ = self.images[i].shape
                result[y:y2 + y][0:x2] = self.images[i]

        return result
