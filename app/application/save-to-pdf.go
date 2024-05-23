package application

import (
	"log"
	"path/filepath"

	"github.com/signintech/gopdf"
	"github.com/vitor-thomazini/video-to-picture/app/domain"
)

type SaveToPdf struct{}

func NewSaveToPdf() *SaveToPdf {
	return &SaveToPdf{}
}

func (p SaveToPdf) Execute(dirPath string, imagesMap map[string][]domain.Resource) {
	log.Println("waiting save pdf")
	var pdf gopdf.GoPdf
	for name, imagesList := range imagesMap {
		log.Printf("creating new pdf with name %s.pdf\n", name)

		pdf.Start(p.config())

		for _, image := range imagesList {
			pdf.AddPage()
			pdf.ImageFrom(image.Image, 0, 0, p.a4PageSize())
		}

		path := filepath.Join(dirPath, name+".whatsapp.pdf")
		pdf.WritePdf(path)
	}
	log.Println("PDF's created with successfully")
}

func (p SaveToPdf) a4PageSize() *gopdf.Rect {
	return &gopdf.Rect{W: 595, H: 842}
}

func (p SaveToPdf) config() gopdf.Config {
	return gopdf.Config{
		PageSize: *p.a4PageSize(),
		Unit:     gopdf.UnitPX,
	}
}
