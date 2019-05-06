package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	bot "github.com/gtaswin/Slackbot/internal"
	"github.com/nlopes/slack"
)

var (
	conf = flag.String("conf", "config/config.ini", "Configuration File")
)

func main() {
	flag.Parse()

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
	cfg, err := ini.Load(*conf)
	if err != nil {
		log.Panic("Fail to read Configuration")
	}
	token := cfg.Section("main").Key("token").String()
	fmt.Println(token)
	debug, _ := strconv.ParseBool(cfg.Section("main").Key("debug").String())
	api := slack.New(token, slack.OptionDebug(debug))
	// api = slack.New(token, slack.OptionDebug(false), slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),)
	rtm := api.NewRTM()
	channels, _ := api.GetChannels(false)

	go rtm.ManageConnection()

	var wg sync.WaitGroup
	for msg := range rtm.IncomingEvents {
		cfg, err := ini.Load(*conf)
		if err != nil {
			log.Error("Fail to read Configuration")
		}
		wg.Add(1)
		go bot.Run(msg, &wg, cfg, channels, rtm)
	}
}
