package dunai

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

func UpdateStars(cv *CV) {
	expr := regexp.MustCompile(`.*github.com/([a-zA-Z0-9_\.-]+)/([a-zA-Z0-9_\.-]+)`)
	star_expr := regexp.MustCompile(`(\d+) users starred`)
	for index, project := range cv.Software {
		matches := expr.FindStringSubmatch(project.Url)
		if len(matches) != 3 {
			continue
		}
		fmt.Println(matches[1], matches[2])
		resp, err := http.Get(fmt.Sprintf("https://github.com/%s/%s/network", matches[1], matches[2]))
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		star_match := star_expr.FindStringSubmatch(string(body))
		if len(star_match) != 2 {
			continue
		}
		stars, _ := strconv.Atoi(star_match[1])
		fmt.Println(stars)
		cv.Software[index].Stars = stars
	}
}
