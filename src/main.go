package main

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://github.com/ethereum/EIPs/tree/master/EIPS"

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("td").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("a").Text()
		href, hasHref := s.Find("a").Attr("href")

		r := regexp.MustCompile(`^eip-.*\.md$`)
		if r.MatchString(text) && hasHref {
			println("====https://github.com" + href + "====")
			scrape("https://github.com" + href)
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
		table.Find("tbody").Find("td").Each(func(idx int, s *goquery.Selection) {
			if len(titleList) > idx {
				title := titleList[idx]
				if title == "eip" || title == "title" || title == "category" || title == "status" {
					println("-> " + title + ": " + s.Text())
				}
			}
		})
	}
}
