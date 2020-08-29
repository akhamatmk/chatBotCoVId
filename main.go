package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	tb "gopkg.in/tucnak/telebot.v2"
)

type global struct {
	NewConfirmed   int64
	TotalConfirmed int64
	NewDeaths      int64
	TotalDeaths    int64
	NewRecovered   int64
	TotalRecovered int64
}

type countrie struct {
	Country        string
	CountryCode    string
	Slug           string
	NewConfirmed   int64
	TotalConfirmed int64
	NewDeaths      int64
	TotalDeaths    int64
	NewRecovered   int64
	TotalRecovered int64
}

type ResponseDataCovid struct {
	Global    global
	Countries []countrie
}

func main() {

	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hi!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		finalMessage := ""
		message := strings.ToUpper(m.Text)
		dataCovid := callCovidData()
		switch message {
		case "CASES TOTAL":
			finalMessage = fmt.Sprintf("Total Active Cases in World %s", humanize.Comma(dataCovid.Global.TotalConfirmed))
			break
		case "DEATHS TOTAL":
			finalMessage = fmt.Sprintf("Total Death Cases in World %s", humanize.Comma(dataCovid.Global.TotalDeaths))
			break
		default:
			finalMessage = otherCommand(message, dataCovid)
		}

		b.Send(m.Sender, finalMessage)
	})

	b.Start()
}

func otherCommand(cmd string, data ResponseDataCovid) string {
	str := strings.Split(cmd, " ")
	if len(str) == 1 {
		return ""
	}

	finalMessage := ""
	switch str[0] {
	case "CASES":

		for i := 0; i < len(data.Countries); i++ {
			if data.Countries[i].CountryCode == strings.ToUpper(str[1]) {
				finalMessage = fmt.Sprintf("Total case in %s =  %s", data.Countries[i].Country, humanize.Comma(data.Countries[i].TotalConfirmed))
				break
			}
		}

		break
	case "DEATHS":
		for i := 0; i < len(data.Countries); i++ {
			if data.Countries[i].CountryCode == strings.ToUpper(str[1]) {
				finalMessage = fmt.Sprintf("In death in %s =  %s", data.Countries[i].Country, humanize.Comma(data.Countries[i].TotalDeaths))
				break
			}
		}

		break
	default:
		finalMessage = "Perintah Orak Ana Jancuk"
	}

	return finalMessage
}

func callCovidData() ResponseDataCovid {
	data := ResponseDataCovid{}
	url := "https://api.covid19api.com/summary"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return data
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return data
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return data
	}

	_ = json.Unmarshal(body, &data)

	return data
}
