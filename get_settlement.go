package klikbca

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type settlementDetail struct {
	Description string `json:"description"`
	Amount      string `json:"amount"`
}

func (klikBca klikBca) GetTodaySettlement() ([]settlementDetail, error) {

	var err error
	settlement := []settlementDetail{}

	klikBca.colly.OnHTML("#pagebody > span > table:last-child > tbody", func(e *colly.HTMLElement) {

		if strings.Contains(string(e.Response.Body), "TIDAK ADA TRANSAKSI") {
			// This should return an empty array
			return
		}

		// Need to skip the first row, because the first row is a header
		// Need to skip the last row, because the last row is the "Saldo awal" dan "Saldo akhir", which we dont need
		tableRow := e.DOM.Children().Slice(1, e.DOM.Children().Length()-1)

		// Need to skip first and last row, because it does not contain anything
		tableData := tableRow.Children().Slice(1, tableRow.Children().Length()-1)

		cleanData := []string{}
		tableData.Find("table > tbody > tr[bgcolor] > td:not([valign])").Each(func(_ int, s *goquery.Selection) {

			// Need to retain the HTML codes, because the data is separated by <br>s which we can extract easily.
			// If we only extract by Text(), we cannot extract the data
			data, err := goquery.OuterHtml(s)
			if err != nil {
				// Just going to panic, because to return this error is kinda hard
				panic(err)
			}

			cleanData = append(cleanData, strings.Split(data, "\n")...)
		})

		// Start processing data
		for _, data := range cleanData {

			// First we sanitize the data by removing </d>
			data = strings.Replace(data, "<td>", "", -1)

			// Second we sanitize the data by remomving </td>
			data = strings.Replace(data, "</td>", "", -1)

			// Then we split the data with <br/> so we can separate the amont data
			splittedData := strings.Split(data, "<br/>")

			// We can rejoin the splitted data to build a description of the transaction
			// Also not forgetting to remove the last item of the array, because it contains the amount
			description := strings.Join(splittedData[0:len(splittedData)-1], " ")

			// The amount should always be in the last item of the array
			amount := splittedData[len(splittedData)-1]

			// Construct a settlement
			settlement = append(settlement, settlementDetail{
				Description: description,
				Amount:      amount,
			})
		}

	})

	// Login
	form := url.Values{}
	form.Add("value(user_id)", klikBca.username)
	form.Add("value(pswd)", klikBca.password)
	form.Add("value(Submit)", "LOGIN")
	form.Add("value(actions)", "login")
	form.Add("value(user_ip)", klikBca.ipAddress)
	form.Add("user_ip", klikBca.ipAddress)
	form.Add("value(mobile)", "true")
	form.Add("mobile", "true")

	err = klikBca.colly.Request(
		"POST",
		"https://m.klikbca.com/authentication.do",
		strings.NewReader(form.Encode()),
		nil,
		http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
	)
	if err != nil {
		return nil, err
	}

	// Go to account information page. This is the place where you can go to the "Balance", "Settlement" etc page
	err = klikBca.colly.Request(
		"POST",
		"https://m.klikbca.com/accountstmt.do?value(actions)=menu",
		nil,
		nil,
		http.Header{},
	)
	if err != nil {
		return nil, err
	}

	// Go to account statement page
	err = klikBca.colly.Request(
		"POST",
		"https://m.klikbca.com/accountstmt.do?value(actions)=acct_stmt",
		nil,
		nil,
		http.Header{},
	)
	if err != nil {
		return nil, err
	}

	// Get formatted time
	now := time.Now()
	date := now.Format("02")
	month := now.Format("01")
	year := now.Format("2006")

	form = url.Values{}
	form.Add("value(D1)", "0")
	form.Add("value(r1)", "1")
	form.Add("value(startDt)", date)
	form.Add("value(startMt)", month)
	form.Add("value(startYr)", year)
	form.Add("value(endDt)", date)
	form.Add("value(endMt)", month)
	form.Add("value(endYr)", year)
	form.Add("value(fDt)", "")
	form.Add("value(tDt)", "")
	form.Add("value(submit1)", "View Account Statement")

	// Get the settlement for today
	err = klikBca.colly.Request(
		"POST",
		"https://m.klikbca.com/accountstmt.do?value(actions)=acctstmtview",
		strings.NewReader(form.Encode()),
		nil,
		http.Header{},
	)
	if err != nil {
		return nil, err
	}

	klikBca.colly.Visit("https://m.klikbca.com/authentication.do")

	return settlement, err
}