import numpy as np
from math import floor
from domain.resource import Resource

from PIL import Image

def draw(image: np.ndarray, resource: Resource, key, count) -> Resource:
    index = floor(resource.counter / 3)

    if (resource.counter % 3 == 0) or (resource.counter == 0):
        resource.add(image)
        # for i in resource.images:
            # debugIm(i, key, count)
    else:
        resource.concat(index, image)
        # for i in resource.images:
        #     debugIm(i, key, count)

    resource.counter += 1
    return resource

def debugIm(img, key, frame):
    image = Image.fromarray(img.astype('uint8'))
    image_path =  "result/%s-%s.png" % (key, str(frame))
    image.save(image_path)

 