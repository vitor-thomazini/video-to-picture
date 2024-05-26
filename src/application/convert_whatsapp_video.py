import cv2 as cv

from application.retrieve_date_from_image import RetrieveDateFromImage
from application.save_pdf import SavePDF
from domain.flow_controller import FlowController
from domain.frame import Frame
from domain.drawer_service import draw
from domain.resource_service import createResource, updateLastResourceByKey

BUFFER_CRITERIA = 15

class ConvertWhatsappVideo:
    def __init__(self, retrieveDateFromImage: RetrieveDateFromImage, savePdf: SavePDF) -> None:
        self.retrieveDateFromImage = retrieveDateFromImage
        self.savePdf = savePdf
        self.ctrl = FlowController()

    def execute(self, filepath, targetDirectory):
        capture = cv.VideoCapture(filepath)
        if not capture.isOpened():
            print("Cannot open camera")
            exit()

        resources = dict()
        keys = set()
        while True:
            ret, frame = capture.read()

            if not ret:
                print("End of stream")
                break
            # self.ctrl.counter +=1
            # print(self.ctrl.counter)
            # if self.ctrl.counter < 175299:
            if self.ctrl.wait():
                continue

            # gray = cv.cvtColor(frame, cv.COLOR_RGB2GRAY)
            texts = self.retrieveDateFromImage.execute(frame)
            frame = cv.resize(frame, (1080, 2340))
            frame = cv.cvtColor(frame, cv.COLOR_BGR2RGB)

            texts.sort()
            for text in texts:
                resources, resource = createResource(text, resources)
                resource = draw(frame, resource, text, self.ctrl.counter)
                resources = updateLastResourceByKey(text, resources, resource)
                keys.add(text)

            print(texts, self.ctrl.text, self.ctrl.counter)
            if len(texts) == 0:
                resources, resource = createResource(self.ctrl.text, resources)
                resource = draw(frame, resource, self.ctrl.text, self.ctrl.counter)
                resources = updateLastResourceByKey(self.ctrl.text, resources, resource)
            else:
                texts.sort(reverse=True)
                self.ctrl.updateText(texts[0])

            if self.__existsAllFramesForKey(list(keys)):
                self.__save(targetDirectory, resources, keys)

        self.__saveAll(targetDirectory, resources, list(keys))       

    def __existsAllFramesForKey(self, keys: list) -> bool:
        keys.sort()
        key = keys[0]
        counter = 0

        for k in list(keys):
            if k != key:
                counter += 1
        return counter > BUFFER_CRITERIA
    
    def __save(self, targetDirectory: str, resources: dict, keys: set):
        ks = list(keys)
        ks.sort()
        key = ks[0]
        self.savePdf.Execute(targetDirectory, key, resources[key])
        keys.remove(key)
        resources.pop(key)
        print(keys)

    def __saveAll(self, targetDirectory: str, resources: dict, keys: list):
        for key in keys:
            self.savePdf.Execute(targetDirectory, key, resources[key])

    

           

            
   
