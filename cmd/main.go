package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	firestore "erc-scraping/internal/firestore"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

// ERCInfo erc info
type ERCInfo struct {
	Eip      string
	Title    string
	Category string
	Status   string
	URL      string
	Type     string
	Created  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	url := os.Getenv("TARGET_URL")

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("td").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("a").Text()
		href, hasHref := s.Find("a").Attr("href")

		r := regexp.MustCompile(`^eip-.*\.md$`)
		if r.MatchString(text) && hasHref {
			domain := os.Getenv("TARGET_DOMAIN")
			scrape(domain + href)
		}
	})
}

func scrape(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	table := doc.Find("table")
	tableType, hasData := table.Attr("data-table-type")
	if tableType == "yaml-metadata" && hasData {
		// get title
		var titleList []string
		table.Find("thead").Find("th").Each(func(idx int, s *goquery.Selection) {
			titleList = append(titleList, s.Text())
		})

		// get data
		ercInfo := make(map[string]string)

		table.Find("tbody").Find("td").Each(func(idx int, s *goquery.Selection) {
			if len(titleList) > idx {
				title := titleList[idx]
				if title == "eip" || title == "title" || title == "category" || title == "status" || title == "type" || title == "created" {
					ercInfo[title] = s.Text()
				}
			}
		})

		// check if the EIP is ERC
		if strings.ToLower(ercInfo["category"]) == "erc" && strings.ToLower(ercInfo["status"]) == "final" {
			ercInfo["url"] = url
			firestore.Save("ercs", ercInfo["eip"], ercInfo)
			log.Print(ercInfo["category"] + ercInfo["eip"])
		}
	}
}
