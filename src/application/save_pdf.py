from fpdf import FPDF
from PIL import Image
import os


class SavePDF:
    def __init__(self) -> None:
        pass

    def Execute(
            self,
            directory: str,
            date: str, 
            resources: list
    ) -> None:

        pdf = FPDF()
        for idx, resource in enumerate(resources):
            pdf.add_page() 
            image = Image.fromarray(resource.fullImage().astype('uint8'), "RGB")
            image_path = "result/images/%s-whatsapp-%s.png" % (date, idx)
            image.save(image_path)
            pdf.image(image_path, 0, 0, 210, 297)
            os.remove(image_path)
        
        print("saving " + date)
        path = "%s/pdf/%s-whatsapp.pdf" % (directory, date)
        pdf.output(path)
        print("saved with successfully " + date)
        