package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type card struct {
	publisher   string
	title       string
	description string
	start       string
}

// var baseURL string = "https://search.naver.com/search.naver?where=news&query=현대카드&sort=1&start=0"
var baseURL string = "https://search.naver.com/"
var keyword string = "현대카드"
var request string = baseURL + "search.naver?where=news&sort=1&query=" + keyword

func main() {
	var results []card
	c := make(chan []card)

	page := getPages(getResponse(request))

	start := 1
	for i := 0; i < page; i++ {
		go getPage(start, c)
		start += 10
	}

	for i := 0; i < page; i++ {
		cards := <-c
		results = append(results, cards...)
	}

	for _, card := range results {
		fmt.Println(card.publisher, card.title, card.start)
	}

}

func getResponse(request string) *http.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(request)
	checkErr(err)
	checkStatus(res)

	return res
}

func getPages(res *http.Response) int {
	defer res.Body.Close()

	result := 0
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".sc_page_inner").Each(func(i int, s *goquery.Selection) {
		result = s.Find("a").Length()
	})

	return result
}

func getPage(start int, mainC chan []card) {
	c := make(chan card)
	var results []card

	fmt.Println("Request URL: ", request+"&start="+strconv.Itoa(start))
	res := getResponse(request + "&start=" + strconv.Itoa(start))
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".news_area")

	searchCards.Each(func(i int, s *goquery.Selection) {
		go createCard(s, c, strconv.Itoa(start))
	})

	for i := 0; i < searchCards.Length(); i++ {
		result := <-c
		results = append(results, result)
	}

	mainC <- results
}

func createCard(s *goquery.Selection, c chan<- card, start string) {
	publisher := s.Find(".info_group > a").First().Text()
	title := s.Find(".news_tit").Text()
	description := s.Find(".dsc_wrap").Text()

	c <- card{publisher: publisher, title: title, description: description, start: start}
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
