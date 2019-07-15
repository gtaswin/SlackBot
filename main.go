package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	Internal "github.com/gtaswin/Slackbot/internal"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		fmt.Println(body)
		w.WriteHeader(http.json({ok: true});)
		time.Sleep(500 * time.Second)

		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: "UXizpRhMmYkaSZjMpqheHLQS"}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(body))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				// log.Info("Received:", ev.Text)
				fmt.Println(eventsAPIEvent.Data)
				fmt.Println(eventsAPIEvent.InnerEvent)
				fmt.Println(ev.Channel)
				reg, _ := regexp.Compile(`\<.*\>`)
				message := reg.ReplaceAllString(ev.Text, "")
				api.PostMessage(ev.Channel, slack.MsgOptionText("GTA", false))
				api.PostMessage(ev.Channel, slack.MsgOptionText(Internal.Format(message, cfg), false))
				fmt.Println("Sent")
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}
