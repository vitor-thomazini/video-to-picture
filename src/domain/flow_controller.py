MAX_VALUE_FRAMES_COUNTER = 25 #18

class FlowController:
    def __init__(self) -> None:
        self.counter = 0
        self.text = ""

    def wait(self) -> bool:
        wait = self.counter%MAX_VALUE_FRAMES_COUNTER != 0
        self.counter += 1
        return wait
    
    def updateText(self, text) -> None:
        self.text = text