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
	queryFormID := `//input[@type="text"][@id="oneid"]`
	queryFormPW := `//input[@type="password"][@id="onepassword"]`
	queryLoginSubmitButton := `//button[@class="btn btn-primary rounded btn-block w-100 z-depth-0 action-login font-weight-bold waves-effect waves-light"]`
	queryPositionSubmitButton := `//button[@id="btnLoginProc"]`

	// 로그인 화면 URL
	currentURL := ""
	if e := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[브라우저] 브라우저 뷰포트 1920x1080으로 맞춤")
			return nil
		}),
		chromedp.Navigate("https://auth.dgist.ac.kr/login/?agentId=19"), // 19=stud.dgist.ac.kr, 22=my.dgist.ac.kr, ...
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] 로그인 페이지 도착. URL: " + currentURL)
			return nil
		}),
		chromedp.Navigate("https://auth.dgist.ac.kr/login/?agentId=19"),
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] 로그인 페이지 재도착. URL: " + currentURL)
			return nil
		}),
		chromedp.WaitVisible(queryFormID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] OneID form 보임")
			return nil
		}),
		chromedp.SendKeys(queryFormID, myID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] OneID 입력")
			return nil
		}),
		chromedp.WaitVisible(queryFormPW),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] 비밀번호 form 보임")
			return nil
		}),
		chromedp.SendKeys(queryFormPW, myPW),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] 비밀번호 입력")
			return nil
		}),
		chromedp.WaitVisible(queryLoginSubmitButton),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] Login 버튼 보임")
			return nil
		}),
		chromedp.Click(queryLoginSubmitButton),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[로그인] Login 버튼 클릭")
			return nil
		}),
		chromedp.Sleep(time.Second*2),
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[신분선택] 신분선택 페이지 도착. URL: " + currentURL)
			targetURL := "https://stud.dgist.ac.kr/sso/agentProc.jsp"
			if targetURL != currentURL {
				return fmt.Errorf("신분선택 URL 이상: %v != %v", targetURL, currentURL)
			}
			return nil
		}),
		chromedp.WaitVisible(queryPositionSubmitButton),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[신분선택] 로그인 버튼 보임")
			return nil
		}),
		chromedp.Click(queryPositionSubmitButton),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[신분선택] 로그인 버튼 클릭")
			return nil
		}),
		/*
			chromedp.ActionFunc(func(ctx context.Context) error {
				log.Println("[로그인][Pop-up] Try to handle pop-up")

				// about:blank 가 열릴 것이다
				newTargetChannel := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
					log.Println("[로그인][Pop-up] chromedp.WaitNewTarget callback")
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
						log.Println("[로그인][Pop-up] Pop-up window URL: " + URL)
						return nil
					}),
					chromedp.WaitVisible(userSelectSubmitButon),
					chromedp.Click(userSelectSubmitButon),
					chromedp.ActionFunc(func(ctx context.Context) error {
						log.Println("[로그인][Pop-up] Click submit button")
						return nil
					}),
				); e != nil {
					return e
				}

				return nil
			}),
		*/
		chromedp.Sleep(time.Second*3),
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[학생정보광장] 로그인 성공. 학생정보광장 URL: " + currentURL)
			targetURL := "https://stud.dgist.ac.kr/sch/student/main.do"
			if targetURL != currentURL {
				return fmt.Errorf("학생정보광장 URL 이상: %v != %v", targetURL, currentURL)
			}
			return nil
		}),
		chromedp.Navigate("https://stud.dgist.ac.kr/usd/usdqSptRechMngtStud/index.do"),
		chromedp.Location(&currentURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[학생정보광장] 전문연구요원>복무상황조회>복무정보 URL: " + currentURL)

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

			log.Println("[학생정보광장] 쿠키 획득 성공")

			return nil
		}),
	); e != nil {
		log.Println(e)
	}

	return cookies
}

func saveCookies(cookies []*http.Cookie) {
	json := `{ "cookies": [] }`

	for _, v := range cookies {
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

func getCookies(doRenew bool) []*http.Cookie {
	stat, e := os.Stat(pathCookieFile)

	if os.IsNotExist(e) || len(gjson.Parse(readFileAsString(pathCookieFile)).Get("cookies").Array()) == 0 {
		log.Println(pathCookieFile + " not exist or empty")

		// cookie file is not exist or empty
		// get cookie from online
		cookies := login(
			config.Get("login.id").String(),
			config.Get("login.pw").String())
		saveCookies(cookies)
		return cookies
	}

	if stat.IsDir() {
		// it's folder
		log.Fatalln(pathCookieFile + " is folder")
		return nil
	}

	// exist but get a renew request
	if doRenew {
		log.Println(pathCookieFile + " exist, remove old cookie")
		if e := os.RemoveAll(pathCookieFile); e != nil {
			log.Fatal(e)
		}

		cookies := login(
			config.Get("login.id").String(),
			config.Get("login.pw").String())
		saveCookies(cookies)
		log.Println(pathCookieFile + " saved")
		return cookies
	}

	// exist
	//log.Println(pathCookieFile + " exist, loading")
	cookies := loadCookies()
	return cookies
}
