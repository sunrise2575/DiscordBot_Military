# DGIST 전문연 복무 알림 디코 봇

DGIST 전문연 복무를 알려주는 디스코드 봇입니다.

출퇴근하거나, 오늘 상태에 업데이트가 있으면, 5초 이내에 DM으로 알려줍니다.

## 사용법

1. `config.example.json` 파일을 수정하고 `config.json` 으로 이름을 바꿉니다. 아래 정보를 기입하면 됩니다.
   - Discord User ID
   - Discord Bot Token
   - stud.dgist.ac.kr 로그인 ID
   - stud.dgist.ac.kr 로그인 PW
   - stud.dgist.ac.kr 학번

2. `go run .` 으로 실행.

## 사용한 라이브러리

   - 디스코드 봇: [discordgo](https://github.com/bwmarrin/discordgo)
   - Chrome headless 로그인 처리: [chromedp](https://github.com/chromedp/chromedp), [cdproto](https://github.com/chromedp/cdproto)
   - JSON 처리: [gjson](https://github.com/tidwall/gjson), [sjson](https://github.com/tidwall/sjson)
   - Cronjob 실행: [cron](https://github.com/robfig/cron)
