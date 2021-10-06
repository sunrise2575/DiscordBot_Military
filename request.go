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
func getMonthInfo(yearAndMonth string) string {
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
	for _, cookie := range myCookies {
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
	for {
		json := gjson.Parse(getMonthInfo(target.Format("200601")))

		/*
			// linear search
			// JSON이 이미 정렬되어 있기에 쓸 필요 없음
			for _, v := range jsonToday.Get("user").Array() {
				candidate, e := time.Parse("2006/01/02", v.Get("WK_APPE_DT").String())
				if candidate.Year() == target.Year() && candidate.Month() == target.Month() && candidate.Day() == target.Day() {
					return v
				}
			}
		*/

		if json.Get("user").Exists() {
			return json.Get("user").Array()[target.Day()-1]
		}

		// 로그인 만료되어서 제대로 메시지가 안오는것임
		log.Println("쿠키 획득 시도")
		myCookies = getCookies()
		log.Println("쿠키 획득 성공")
	}
}
