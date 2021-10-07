package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func login(myID, myPW string) (cookies []*http.Cookie) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	// 로그인 화면 DOM path
	formID := `//input[@id="username"][@class="form-control"]`
	formPW := `//input[@id="passwordTest"][@class="form-control"]`
	formSubmitButton := `//input[@id="btnLogin"][@class="btn btn-info login_btn"]`

	// 로그인 화면 URL
	loginPageURL := ""
	if e := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Try to login")
			return nil
		}),
		chromedp.EmulateViewport(1920, 1080),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Set browser viewport to 1920x1080")
			return nil
		}),
		chromedp.Navigate("https://stud.dgist.ac.kr/login.jsp"),
		chromedp.Location(&loginPageURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Log-in page URL: " + loginPageURL)
			return nil
		}),
		chromedp.WaitVisible(formID),
		chromedp.SendKeys(formID, myID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Input ID")
			return nil
		}),
		chromedp.WaitVisible(formPW),
		chromedp.SendKeys(formPW, myPW),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Input password")
			return nil
		}),
		chromedp.WaitVisible(formSubmitButton),
		chromedp.Click(formSubmitButton),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login] Click submit button")
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login][Pop-up] Try to handle pop-up")

			// about:blank 가 열릴 것이다
			newTargetChannel := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
				return info.URL != ""
			})

			blankCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(<-newTargetChannel))
			defer cancel()

			// 신분 및 보직 선택
			URL := ""
			userSelectSubmitButon := `//button[@id="btnSelect"][@class="btn_select"]`
			if e := chromedp.Run(blankCtx,
				chromedp.Location(&URL),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("[Login][Pop-up] Pop-up window URL: " + URL)
					return nil
				}),
				chromedp.WaitVisible(userSelectSubmitButon),
				chromedp.Click(userSelectSubmitButon),
				chromedp.ActionFunc(func(ctx context.Context) error {
					log.Println("[Login][Pop-up] Click submit button")
					return nil
				}),
			); e != nil {
				return e
			}

			return nil
		}),
		chromedp.Sleep(time.Second*2),
		chromedp.Location(&loginPageURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Login Complete] New page URL: " + loginPageURL)
			targetURL := "https://stud.dgist.ac.kr/sch/student/main.do"
			if targetURL != loginPageURL {
				log.Println("[Login Failed?] Wrong target!")
			}
			return nil
		}),
		chromedp.Navigate("https://stud.dgist.ac.kr/usd/usdqSptRechMngtStud/index.do"),
		chromedp.Location(&loginPageURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[Cookie] New page URL: " + loginPageURL)

			_cookiesTemp, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, cookie := range _cookiesTemp {
				if cookie.Domain == "stud.dgist.ac.kr" || cookie.Domain == ".dgist.ac.kr" {
					// 필요한 쿠키만 결과에 삽입
					cookies = append(cookies, &http.Cookie{
						Name:  cookie.Name,
						Value: cookie.Value,
					})
				}
			}

			log.Println("[Cookie] Complete getting the cookie!")

			return nil
		}),
	); e != nil {
		log.Fatal(e)
	}

	return cookies
}

func saveCookies() {
	json := `{ "cookies": [] }`

	for _, v := range myCookies {
		json, _ = sjson.SetRaw(json, "cookies.-1", fmt.Sprintf(`{"Name":"%v", "Value":"%v"}`, v.Name, v.Value))
	}

	if e := ioutil.WriteFile(pathCookieFile, []byte(json), 0644); e != nil {
		log.Println(e)
	}
}

func loadCookies() []*http.Cookie {
	cookies := []*http.Cookie{}

	for _, v := range gjson.Parse(readFileAsString(pathCookieFile)).Get("cookies").Array() {
		cookies = append(cookies, &http.Cookie{
			Name:  v.Get("Name").String(),
			Value: v.Get("Value").String(),
		})
	}

	return cookies
}

func getCookies() []*http.Cookie {
	stat, e := os.Stat(pathCookieFile)

	if os.IsNotExist(e) || len(gjson.Parse(readFileAsString(pathCookieFile)).Get("cookies").Array()) == 0 {
		log.Println(pathCookieFile + " not exist or empty")

		// cookie file is not exist or empty
		// get cookie from online
		cookies := login(
			config.Get("login.id").String(),
			config.Get("login.pw").String())
		myCookies = cookies
		saveCookies()
		return cookies
	}

	if stat.IsDir() {
		// it's folder
		log.Fatalln(pathCookieFile + " is folder")
		return nil
	}

	// exist
	log.Println(pathCookieFile + " exist")
	cookies := loadCookies()
	return cookies

}
