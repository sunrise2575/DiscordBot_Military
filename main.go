package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

var (
	config         gjson.Result
	pathCookieFile = "./cookie.json"
	pathConfigFile = "./config.json"
	pathTodayFile  = "./today.json"
	studentID      = ""
)

func init() {
	config = gjson.Parse(readFileAsString(pathConfigFile))
	studentID = config.Get("login.studentID").String()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create discord session
	discord, e := discordgo.New("Bot " + config.Get("discord.botToken").String())
	if e != nil {
		log.Fatalln("error creating Discord session,", e)
		return
	}

	channel, e := discord.UserChannelCreate(config.Get("discord.userID").String())
	if e != nil {
		log.Println("unable to create user channel", e)
	}

	// 크론잡 등록
	c := cron.New(cron.WithSeconds())
	c.Start()
	c.AddFunc("*/15 * * * * *", func() {
		checkStatusCron(discord, channel)
	})

	// 메인 함수가 종료되면 실행될 것들
	defer func() {
		c.Stop()
		discord.Close()
		log.Println("bye")
	}()

	// 세팅을 다 했으니 세션을 연다
	if e := discord.Open(); e != nil {
		log.Fatalln("error opening connection,", e)
		return
	}

	// 봇의 상태 변경
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

	// Ctrl+C를 받아서 프로그램 자체를 종료하는 부분. os 신호를 받는다
	log.Println("bot is now running. Press Ctrl+C to exit.")
	{
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}
	log.Println("received Ctrl+C, please wait.")
}
