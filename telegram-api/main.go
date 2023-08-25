package telegramapi

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"

	config "github.com/southouse/naver-news-crawler/conf"
)

const telegramAPIUrl string = "https://api.telegram.org/bot"
const path string = "/sendmessage"

func SendMessage(text string) {
	config, err := config.GetConfig("config.yaml")
	checkErr(err)

	requestUrl := telegramAPIUrl + config.Telegram.ApiKey + "/" + path + "?chat_id=" + config.Telegram.ChatId + "&text=" + text
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", requestUrl, nil)
	// fmt.Println(requestUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	checkErr(err)

	res, err := client.Do(req)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

	_, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalln(err)
	}

	fmt.Println(requestUrl)
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
