package application

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type GetTextFromImage struct {
	client          *gosseract.Client
	allTextMatchers []string
}

func NewGetTextFromImage() GetTextFromImage {
	return GetTextFromImage{
		client: gosseract.NewClient(),
		allTextMatchers: []string{
			`.*(Sunday).*`, `.*(Monday).*`, `.*(Tuesday).*`, `.*(Wednesday).*`, `.*(Thursday).*`, `.*(Friday).*`, `.*(Saturday).*`,
			`.*(Yesterday).*`, `.*(Today).*`, `.*(January \d{2},\d{4}).*`, `.*(February \d{2},\d{4}).*`, `.*(March \d{2},\d{4}).*`,
			`.*(April \d{2},\d{4}).*`, `.*(May \d{2},\d{4}).*`, `.*(June \d{2},\d{4}).*`, `.*(July \d{2},\d{4}).*`, `.*(August \d{2},\d{4}).*`,
			`.*(September \d{2},\d{4}).*`, `.*(October \d{2},\d{4}).*`, `.*(November \d{2},\d{4}).*`, `.*(December \d{2},\d{4}).*`,
		},
	}
}

func (g GetTextFromImage) Execute(image image.Image) []string {
	buff := new(bytes.Buffer)
	err := png.Encode(buff, image)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	g.client.SetImageFromBytes(buff.Bytes())
	text, _ := g.client.Text()
	// Monday ini May 3, 2024
	phrases := strings.Split(text, "\n")

	labels := make([]string, 0)
	for _, phrase := range phrases {
		for _, regex := range g.allTextMatchers {
			r := regexp.MustCompile(regex)
			if r.Match([]byte(phrase)) {
				label := r.FindStringSubmatch(phrase)[1]
				labels = append(labels, label)
			}
		}

	}

	return labels
}

func (g GetTextFromImage) CloseTransaction() {
	g.client.Close()
}
