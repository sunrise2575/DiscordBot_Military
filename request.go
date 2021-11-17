package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// the format of yearAndMonth is "200601"
func getMonthInfo(cookies []*http.Cookie, yearAndMonth string) string {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("searchStdNo", studentID)
	_ = writer.WriteField("searchYymm", yearAndMonth)

	err := writer.Close()
	if err != nil {
		log.Println(err)
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://stud.dgist.ac.kr/usd/usdqSptRechMngtStud/listPbsvAppe.do", payload)
	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	return string(body)
}

func getDayInfo(target time.Time) gjson.Result {
	cookies := []*http.Cookie{}

	for {
		//log.Println("cookies:", cookies)
		// load cookie if not exist
		if len(cookies) == 0 {
			//log.Println("getCookies(false)")
			cookies = getCookies(false)
			//log.Println("cookies:", cookies)
		}

		if len(cookies) > 0 {
			//log.Println("getMonthInfo")
			for i := 0; i < 5; i++ {
				json := gjson.Parse(getMonthInfo(cookies, target.Format("200601")))

				//log.Println("cookies:", cookies)
				//log.Println("json:", json.String())
				if json.Get("user").Exists() {
					//log.Println("yes")
					// operating normally
					return json.Get("user").Array()[target.Day()-1]
				}
				time.Sleep(time.Millisecond * 100)
			}
		}

		// somethings wrong
		// renew cookie
		log.Println("cookie may be expired, try to renew cookie...")
		cookies = getCookies(true)
		if cookies != nil {
			log.Println("success to renew cookie")
			continue
		}

		log.Println("failed, try more..")

		// failed. try more...
		time.Sleep(time.Second * 1)
	}
}
