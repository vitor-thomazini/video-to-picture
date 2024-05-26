from application.convert_whatsapp_video import ConvertWhatsappVideo
from application.retrieve_date_from_image import RetrieveDateFromImage
from application.save_pdf import SavePDF

s = SavePDF()
r = RetrieveDateFromImage()
c = ConvertWhatsappVideo(r, s)
c.execute("/Users/Vitor/Desktop/video.mp4", "result")

