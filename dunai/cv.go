package dunai

import (
	"fmt"
	// "github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
	"time"
)

type Company struct {
	Name         string              `yaml:"name"`
	Icon         string              `yaml:"icon"`
	Color        string              `yaml:"color"`
	Role         string              `yaml:"role"`
	Technologies []map[string]string `yaml:"technologies"`
	Info         string              `yaml:"info"`
	Description  []string            `yaml:"description"`
	Start        *time.Time          `yaml:"start"`
	End          *time.Time          `yaml:"end"`
}

type Software struct {
	Name        string `yaml:"name"`
	Url         string `yaml:"url"`
	Description string `yaml:"description"`
	Stars       int    `yaml:"stars"`
}

type Project struct {
	Name        string `yaml:"name"`
	Url         string `yaml:"url"`
	Description string `yaml:"description"`
	Image       string `yaml:"image"`
}

type Contact struct {
	Title string `yaml:"title"`
	Url   string `yaml:"url"`
}

type CV struct {
	Bio       []string   `yaml:"bio"`
	Features  []string   `yaml:"features"`
	Education []string   `yaml:"education"`
	Languages []string   `yaml:"languages"`
	Companies []Company  `yaml:"companies"`
	Software  []Software `yaml:"software"`
	Projects  []Project  `yaml:"projects"`
	Misc      []string
	Contacts  []Contact
}

// Map of month names to Ukrainian
var monthNames = map[string]string{
	"January":   "січень",
	"February":  "лютий",
	"March":     "березень",
	"April":     "квітень",
	"May":       "травень",
	"June":      "червень",
	"July":      "липень",
	"August":    "серпень",
	"September": "вересень",
	"October":   "жовтень",
	"November":  "листопад",
	"December":  "грудень",
}

func (c Company) PrettyStart() string {
	return fmt.Sprintf("%s %d", monthNames[c.Start.Month().String()], c.Start.Year())
}

func (c Company) PrettyEnd() string {
	return fmt.Sprintf("%s %d", monthNames[c.End.Month().String()], c.End.Year())
}

func ReadCV() (*CV, error) {
	cv := &CV{}
	// data, err := ioutil.ReadFile("./data/cv.yaml")
	// if err != nil {
	// 	return nil, err
	// }
	if err := yaml.Unmarshal(cvYaml, cv); err != nil {
		return nil, err
	}
	return cv, nil
}

// func (cv *CV) GeneratePDF() (*gopdf.GoPdf, error) {
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
// 	pdf.AddPage()
// 	err := pdf.AddTTFFont("dejavu", "./data/DejaVuSans.ttf")
// 	if err != nil {
// 		return nil, err
// 	}
// 	pdf.SetFont("dejavu", "", 14)
// 	pdf.Cell(nil, "TODO")
// 	return &pdf, nil
// }
