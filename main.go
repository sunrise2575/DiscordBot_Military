package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/tidwall/gjson"
)

func readFileAsString(path string) string {
	out, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	return string(out)
}

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

// the format of yearAndMonth is "200601"
func getDailyInfo(cookies []*http.Cookie, studentID, yearAndMonth string) string {
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	/*
		// create discord session
		discord, e := discordgo.New("Bot " + readFileAsString("./config/token.txt"))
		if e != nil {
			log.Fatalln("error creating Discord session,", e)
			return
		}

		// 세팅을 다 했으니 세션을 연다
		if e := discord.Open(); e != nil {
			log.Fatalln("error opening connection,", e)
			return
		}

		// 메인 함수가 종료되면 실행될 것들
		defer func() {
			discord.Close()
			log.Println("bye")
		}()

		if e := discord.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "육군 아미 타이거",
					Type: discordgo.ActivityTypeGame,
				},
			},
		}); e != nil {
			log.Fatalln("error update status complex,", e)
			return
		}
	*/

	// login
	config := gjson.Parse(readFileAsString("./config.json"))
	cookies := login(
		config.Get("login.id").String(),
		config.Get("login.pw").String())
	// save cookie to db
	now := time.Now()
	json := gjson.Parse(getDailyInfo(cookies,
		config.Get("login.studentID").String(),
		now.Format("200601")))

	for _, v := range json.Get("user").Array() {
		target, e := time.Parse("2006/01/02", v.Get("WK_APPE_DT").String())
		if e != nil {
			log.Println(e)
		}

		if target.Year() == now.Year() &&
			target.Month() == now.Month() &&
			target.Day() == now.Day() {
			if v.Get("FROM_HM").Exists() {
				log.Printf("오늘 당신은 %v시에 출근했습니다.\n", v.Get("FROM_HM"))
			}
			if v.Get("TO_HM").Exists() {
				log.Printf("오늘 당신은 %v시에 퇴근했습니다.\n", v.Get("TO_HM"))
			}
		}

	}

	/*
			channel, e := discord.UserChannelCreate(readFileAsString("./config/discord_user_id.txt"))
			if e != nil {
				log.Println("unable to create user channel", e)
			}
			discord.ChannelMessageSend(channel.ID, "")

		// Ctrl+C를 받아서 프로그램 자체를 종료하는 부분. os 신호를 받는다
		log.Println("bot is now running. Press Ctrl+C to exit.")
		{
			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
			<-sc
		}
		log.Println("received Ctrl+C, please wait.")
	*/
}
