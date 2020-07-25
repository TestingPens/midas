package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func getchart(symbol string, date time.Time, resolution string, token string) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, time.UTC)
	fromunix := strconv.FormatInt(from.Unix(), 10)
	to := time.Date(date.Year(), date.Month(), date.Day()+1, 0, 0, 0, 0, time.UTC)
	tounix := strconv.FormatInt(to.Unix(), 10)
	url := fmt.Sprintf("https://apidojo-yahoo-finance-v1.p.rapidapi.com/stock/v2/get-chart?interval=5m&region=US&symbol=%s&lang=en&period1=%s&period2=%s", symbol, fromunix, tounix)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Host", "apidojo-yahoo-finance-v1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	var chart Chart

	if resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&chart)
		if err != nil {
			panic(err)
		}
		writechart(&chart)
	} else {
		fmt.Println("Error retrieving: " + symbol)
	}
	if chart.Chart.Error != nil {
		fmt.Println("Error retrieving: " + symbol)
	}
}
