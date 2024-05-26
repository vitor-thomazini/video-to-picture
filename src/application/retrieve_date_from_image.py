import re
import pytesseract
import numpy

from datetime import datetime

LATEST_DAY_REGEX = "(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Yesterday|Today)"
DAY_REGEX        = "(January \d{1,2}, \d{4}|February \d{1,2}, \d{4}|March \d{1,2}, \d{4}|April \d{1,2}, \d{4}|May \d{1,2}, \d{4}|June \d{1,2}, \d{4}|July \d{1,2}, \d{4}|August \d{1,2}, \d{4}|September \d{1,2}, \d{4}|October \d{1,2}, \d{4}|November \d{1,2}, \d{4}|December \d{1,2}, \d{4})"

class RetrieveDateFromImage:
    def __init__(self) -> None:
        pass

    def execute(self, image: numpy.ndarray) -> list:
        fullText = pytesseract.image_to_string(image, lang="eng")
        # fullText = fullText.replace(" ","")
        fullText = fullText.replace("\n","")

        texts = re.findall(DAY_REGEX, fullText)
        texts = [*texts, *["latest" for i in re.findall(LATEST_DAY_REGEX, fullText)]]
        texts = set(texts)
 
        texts = [self._convertDate(text) for text in texts]

        return texts

    def _convertDate(self, text: set) -> str:
        if text == "latest":
            return text
        
        date_object = datetime.strptime(text, "%B %d, %Y")
        return date_object.strftime("%Y%m%d")
