package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func checkStatusCron(sess *discordgo.Session, channel *discordgo.Channel) {
	now := time.Now()

	today := getDayInfo(now)
	if existAndIsFile(pathTodayFile) {
		today_old := readFileAsString(pathTodayFile)
		if today_old == today.String() {
			// 오래된것과 비교해서 달라진게 없으면 나간다
			return
		}
	}
	// 뭔가 달라진걸 감지했다. 일단 파일을 기록한다.
	writeFileFromString(pathTodayFile, today.String())

	tomorrow := getDayInfo(now.Add(time.Hour * 24))

	prefix := fmt.Sprintf("[%v]", now.Format("2006-01-02 15:04:05"))
	msg := ""

	if today.Get("FROM_HM").Exists() {
		msg += fmt.Sprintf("출근 감지 (%v)", today.Get("FROM_HM"))
	}

	if today.Get("TO_HM").Exists() {
		if msg != "" {
			msg += ", "
		}
		msg += fmt.Sprintf("퇴근 감지 (%v) %v", today.Get("TO_HM"), today.Get("REMK"))
	}

	if msg == "" {
		// 출근도 퇴근도 없는데 (msg=="") 이 라인에 진입했다면?
		// 당일 휴가를 썼다는 것이다
		msg += fmt.Sprintf("오늘 정보 변경 감지 %v %v %v", today.Get("WK_SCH"), today.Get("ACCP_YN"), today.Get("REMK"))
	}

	suffix := ""

	if tomorrow.Get("WK_SCH").String() == "평상근무" {
		// 내일이 평상근무라면
		if today.Get("WK_SCH").String() == "평상근무" {
			// 오늘도 평상근무라면
			suffix += "내일도 출근"
		} else {
			// 오늘은 평상근무가 아니라면
			suffix += "내일은 출근"
		}
	} else {
		// 내일이 평상근무가 아니라면
		suffix += fmt.Sprintf("내일은 %v %v", tomorrow.Get("WK_SCH"), tomorrow.Get("REMK"))
	}

	log.Println(prefix + " " + msg + ", " + suffix)
	sess.ChannelMessageSend(channel.ID, prefix+" "+msg+", "+suffix)
}
