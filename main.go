package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	"github.com/gtaswin/Slackbot/internal"
	"github.com/nlopes/slack"
)

var (
	conf = flag.String("conf", "config/config.ini", "Configuration File")
)

func main() {
	flag.Parse()

	//Logging setup
	Logfilename := "bot.log"
	f, err := os.OpenFile(Logfilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}

	//Loading Configuration
	cfg, err := ini.Load(*conf)
	if err != nil {
		log.Panic("Fail to read Configuration")
	}

	//Vars for Slack connection to initialise
	token := cfg.Section("main").Key("token").String()
	debug, _ := strconv.ParseBool(cfg.Section("main").Key("debug").String())
	api := slack.New(token, slack.OptionDebug(debug))
	// api = slack.New(token, slack.OptionDebug(false), slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),)
	rtm := api.NewRTM()
	channels, _ := api.GetChannels(false)

	//Cron job function
	go func() {
		Cronchan := make(chan string)
		internal.Cron(cfg, Cronchan)
		for {
			rtm.SendMessage(rtm.NewOutgoingMessage(<-Cronchan, cfg.Section("main").Key("cron_channel").String()))
		}
	}()

	//RTM Management in Slack
	go rtm.ManageConnection()
	var wg sync.WaitGroup
	for msg := range rtm.IncomingEvents {
		cfg, err := ini.Load(*conf)
		if err != nil {
			log.Error("Fail to read Configuration")
		}
		wg.Add(1)
		go internal.Run(msg, &wg, cfg, channels, rtm)
	}
}
