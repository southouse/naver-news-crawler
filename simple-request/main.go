package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// var baseURL string = "https://search.naver.com/search.naver?where=news&query=현대카드&sort=1&start=0"
var baseURL string = "https://search.naver.com/"
var keyword string = "현대카드"
var requestURL string = baseURL + "search.naver?where=news&sort=1&query=" + keyword

func main() {

	getPage()

}

func getResponse(requestURL string) *http.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	checkErr(err)

	res, err := client.Do(req)
	checkErr(err)
	checkStatus(res)

	return res
}

func getPageCount(res *http.Response) int {
	defer res.Body.Close()

	result := 0
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".sc_page_inner").Each(func(i int, s *goquery.Selection) {
		result = s.Find("a").Length()
	})

	return result
}

func getPage() {
	fmt.Println("Request URL: ", requestURL+"&start=1")
	res := getResponse(requestURL + "&start=1")
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".news_area")

	searchCards.Each(func(i int, s *goquery.Selection) {
		publisher := s.Find(".info_group > a").First().Text()
		title := s.Find(".news_tit").Text()
		description := s.Find(".dsc_wrap").Text()

		fmt.Printf("publisher: "+publisher, "title: "+title, "description: "+description)
	})
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("caused error")
		log.Fatalln(err)
	}
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
