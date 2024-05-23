package application

import (
	"log"
	"path/filepath"

	"github.com/signintech/gopdf"
	"github.com/vitor-thomazini/video-to-picture/app/domain"
)

type SaveOneDateInPdf struct{}

func NewSaveOneDateInPdf() *SaveOneDateInPdf {
	return &SaveOneDateInPdf{}
}

func (p SaveOneDateInPdf) Execute(dirPath string, date string, imagesMap []domain.Resource) {
	log.Println("waiting save pdf")
	var pdf gopdf.GoPdf
	log.Printf("creating new pdf with name %s.pdf\n", date)
	defer pdf.Close()

	pdf.Start(p.config())

	for _, resource := range imagesMap {
		pdf.AddPage()
		pdf.ImageFrom(resource.Image(), 0, 0, p.a4PageSize())

		path := filepath.Join(dirPath, date+".whatsapp.pdf")
		pdf.WritePdf(path)
	}

	log.Println("PDF's created with successfully")
}

func (p SaveOneDateInPdf) a4PageSize() *gopdf.Rect {
	return &gopdf.Rect{W: 595, H: 842}
}

func (p SaveOneDateInPdf) config() gopdf.Config {
	return gopdf.Config{
		PageSize: *p.a4PageSize(),
		Unit:     gopdf.UnitPX,
	}
}
