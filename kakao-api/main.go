package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Message struct {
	ObjectType  string `json:"object_type,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        Link   `json:"link,omitempty"`
	ButtonTitle string `json:"buttion_title,omitempty"`
}

type Link struct {
	WebUrl       string `json:"web_url,omitempty"`
	MobileWebUrl string `json:"mobile_web_url,omitempty"`
}

const codeURL string = ""
const oauthURL string = ""
const redirectURL string = ""
const messageURL string = ""
const code string = ""
const client_id string = ""

func getCode() string {
	requestUrl := codeURL + "?client_id=" + client_id + "&redirect_uri=" + redirectURL + "&response_type=code"
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
	return data
}

func getAccessToken() string {
	client := &http.Client{}

	bodyData := url.Values{
		"grant_type":   {"authorization_code"},
		"client_id":    {client_id},
		"redirect_uri": {redirectURL},
		"code":         {code},
	}

	req, err := http.NewRequest("POST", oauthURL, bytes.NewBufferString(bodyData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	return data
}

func sendMessage() {
	accessToken := getAccessToken()
	fmt.Println("accessToken: " + accessToken)

	client := &http.Client{}

	reqBody := Message{
		ObjectType: "text",
		Text:       "테스트입니다.",
		Link: Link{
			WebUrl:       "https://developers.kakao.com",
			MobileWebUrl: "https://developers.kakao.com",
		},
		ButtonTitle: "키워드",
	}

	jsonBody, err := json.Marshal(reqBody)
	checkErr(err)

	req, err := http.NewRequest("POST", messageURL, strings.NewReader(string(jsonBody)))
	req.Header.Set("Authorization", "Authorization: Bearer "+accessToken)
	checkErr(err)

	res, err := client.Do(req)
	checkErr(err)
	checkStatus(res)

	defer res.Body.Close()

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

func main() {
	getAccessToken()
}
