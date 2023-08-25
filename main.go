package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	telegramapi "github.com/southouse/naver-news-crawler/telegram-api"
)

type card struct {
	publisher   string
	title       string
	description string
	url         string
	date        time.Time
}

// var baseURL string = "https://search.naver.com/search.naver?where=news&query=현대카드&sort=1&start=0"
var baseURL string = "https://search.naver.com/"
var keyword string = "%EB%8D%B0%EC%9D%B4%EC%8B%9D%EC%8A%A4%20%22%EC%BD%98%EC%84%9C%ED%8A%B8%22%2B%EA%B0%9C%EC%B5%9C%20-%EC%9D%8C%EC%9B%90%20-%EC%95%A8%EB%B2%94"
var requestURL string = baseURL + "search.naver?where=news&sort=1&query=" + keyword

func main() {
	c := make(chan card)

	pageCount := getPageCount(getResponse(requestURL)) / 2

	start := 1

	for i := 0; i < pageCount; i++ {
		go getPage(start, c)
		start += 10
	}

	for card := range c {
		request := "[" + card.date.Format(time.RFC3339) + "]%0A" + card.publisher + "%0A" + card.title + "%0A" + card.url + "%0A%0A%0A"
		telegramapi.SendMessage(request)
		return
	}
}

func getResponse(requestURL string) *http.Response {
	time.Sleep(1000)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", requestURL, nil)
	// req.Header.Set("User-Agent", "Mozilla/5.0")
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

func getPage(start int, mainC chan<- card) {
	c := make(chan card)

	// fmt.Println("Request URL: ", requestURL+"&start="+strconv.Itoa(start))
	res := getResponse(requestURL + "&start=" + strconv.Itoa(start))
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".news_area")

	searchCards.Each(func(i int, s *goquery.Selection) {
		go createCard(s, c, strconv.Itoa(start))
	})

	for i := range c {
		if IsWithinRange(i.date, time.Now().AddDate(0, 0, -2), time.Now()) {
			mainC <- i
		}
	}

	close(mainC)
}

func createCard(s *goquery.Selection, c chan<- card, start string) {
	var cardSet card

	publisher := s.Find(".info_group > a").First().Text()
	title := s.Find(".news_tit").Text()
	description := s.Find(".dsc_wrap").Text()
	dateString := s.Find(".info_group > span").Text()
	url, _ := s.Find(".news_tit").Attr("href")
	pieces := strings.Split(dateString, ".")

	year := 0
	month := 0
	day := 0
	date := time.Now()

	if !strings.Contains(pieces[0], "시간") && !strings.Contains(pieces[0], "일") {
		year, _ = strconv.Atoi(pieces[0])
		month, _ = strconv.Atoi(pieces[1])
		day, _ = strconv.Atoi(pieces[2])
		date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())
	}

	cardSet.publisher = publisher
	cardSet.title = title
	cardSet.description = description
	cardSet.date = date
	cardSet.url = url

	c <- cardSet
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("caused error")
		log.Fatalln(err)
	}
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(b))
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func IsWithinRange(checkTime, startTime, endTime time.Time) bool {
	return checkTime.After(startTime) && checkTime.Before(endTime)
}
