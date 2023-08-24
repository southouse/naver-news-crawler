package telegramapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const telegramAPIUrl string = "https://api.telegram.org/bot"
const apiKey string = ""
const path string = "/sendmessage"
const chatId string = ""

func SendMessage(text string) {
	requestUrl := telegramAPIUrl + apiKey + "/sendmessage" + "?chat_id=" + chatId + "&text=" + text
	client := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)
	fmt.Println(requestUrl)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	checkErr(err)

	res, err := client.Do(req)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	data := string(body)
	fmt.Println(data)
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
