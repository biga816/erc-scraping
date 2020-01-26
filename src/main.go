package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		ercInfo := ERCInfo{}
		p := reflect.ValueOf(&ercInfo).Elem()
		table.Find("tbody").Find("td").Each(func(idx int, s *goquery.Selection) {
			if len(titleList) > idx {
				title := titleList[idx]
				if title == "eip" || title == "title" || title == "category" || title == "status" || title == "type" || title == "created" {
					// println("-> " + title + ": " + s.Text())
					fieldName := strings.ToUpper(title[0:1]) + title[1:len(title)]
					p.FieldByName(fieldName).SetString(s.Text())
				}
			}
		})

		// check if the EIP is ERC
		if strings.ToLower(ercInfo.Category) == "erc" && strings.ToLower(ercInfo.Status) == "final" {
			ercInfo.URL = url
			fmt.Printf("%+v\n", ercInfo)
		}
	}
}
