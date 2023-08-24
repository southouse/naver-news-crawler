package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var globalValue int

func main() {
	layout := "1970.01.01."
	date := "16시간 전"

	pieces := strings.Split(date, ".")
	fmt.Println(pieces)

	if strings.Contains(pieces[0], "시간") {
		pieces[0] = strconv.FormatInt(int64(time.Now().Year()), 10)
	}

	year, _ := strconv.Atoi(pieces[0])
	month, _ := strconv.Atoi(pieces[1])
	day, _ := strconv.Atoi(pieces[2])

	abc := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())

	time, _ := time.Parse(date, layout)
	fmt.Println(time)
	fmt.Println(abc)
}

func action(i int) {
	fmt.Print(i, " ")

	time.Sleep(time.Second)
}
