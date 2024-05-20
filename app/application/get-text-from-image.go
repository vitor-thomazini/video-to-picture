package application

import (
	"bytes"
	"fmt"
	"image/png"
	"regexp"
	"sort"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

const (
	LATEST_DAY_REGEX = `(Sunday|Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Yesterday|Today)`
	DAY_REGEX        = `(January\d{1,2},\d{4}|February\d{1,2},\d{4}|March\d{1,2},\d{4}|April\d{1,2},\d{4}|May\d{1,2},\d{4}|June\d{1,2},\d{4|July\d{1,2},\d{4}|August\d{1,2},\d{4}|September\d{1,2},\d{4}|October\d{1,2},\d{4}|November\d{1,2},\d{4}|December\d{1,2},\d{4})`
)

type GetTextFromImage struct {
	client          *gosseract.Client
	allTextMatchers []string
}

// Pesquisa sobre May 3, 2024 dus
//

func NewGetTextFromImage() GetTextFromImage {
	return GetTextFromImage{
		client: gosseract.NewClient(),
	}
}

func (g GetTextFromImage) Execute(img gocv.Mat) []string {
	month := map[string]string{
		"January":   "01",
		"February":  "02",
		"March":     "03",
		"April":     "04",
		"May":       "05",
		"June":      "06",
		"July":      "07",
		"August":    "08",
		"September": "09",
		"October":   "10",
		"November":  "11",
		"December":  "12",
	}

	detectionImg := gocv.NewMat()
	gocv.CvtColor(img, &detectionImg, gocv.ColorBGRToGray)

	image, _ := detectionImg.ToImage()

	g.client.SetLanguage("eng", "por")
	buff := new(bytes.Buffer)
	err := png.Encode(buff, image)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	g.client.SetImageFromBytes(buff.Bytes())
	text, _ := g.client.Text()
	text = strings.ReplaceAll(text, " ", "")

	labels := make(map[string]int, 0)

	r := regexp.MustCompile(LATEST_DAY_REGEX)
	if r.Match([]byte(text)) {
		for range r.FindAllStringSubmatch(text, 100) {
			labels["latest"] += 1
		}
	}

	r = regexp.MustCompile(DAY_REGEX)
	if r.Match([]byte(text)) {
		for _, value := range r.FindAllStringSubmatch(text, 100) {
			rd := regexp.MustCompile(`([a-zA-Z]*)([0-9]*),([0-9]*)`)
			v := rd.FindAllStringSubmatch(value[1], 3)
			key := fmt.Sprintf("%4s-%02s-%02s", v[0][3], month[v[0][1]], v[0][2])
			labels[key] += 1
		}
	}

	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func (g GetTextFromImage) CloseTransaction() {
	g.client.Close()
}
