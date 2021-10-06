# discord.bot.notify_military

DGIST 전문연 복무를 알려주는 디스코드 봇입니다.
사용한 라이브러리:
   - [discordgo](https://github.com/bwmarrin/discordgo)
   - [chromedp](https://github.com/chromedp/chromedp), [cdproto](https://github.com/chromedp/cdproto)
   - JSON 처리: [gjson](https://github.com/tidwall/gjson), [sjson](https://github.com/tidwall/sjson)
   - [cron](https://github.com/robfig/cron)

## 사용법

1. `config.example.json` 파일을 수정하고 `config.json` 으로 이름을 바꿉니다. 아래 정보를 기입하면 됩니다.
   - Discord User ID
   - Discord Bot Token
   - stud.dgist.ac.kr 로그인 ID
   - stud.dgist.ac.kr 로그인 PW
   - stud.dgist.ac.kr 학번

2. `go run .` 으로 실행하시면 됩니다.