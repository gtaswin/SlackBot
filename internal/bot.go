package internal

import (
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	"github.com/nlopes/slack"
)

//Run for running according to RTM events
func Run(msg slack.RTMEvent, wg *sync.WaitGroup, cfg *ini.File, channels []slack.Channel, rtm *slack.RTM) {
Switching:
	switch ev := msg.Data.(type) {
	// case *slack.HelloEvent:
	// Ignore hello

	// case *slack.ConnectedEvent:
	// 	for _, channel := range channels {
	// 		rtm.SendMessage(rtm.NewOutgoingMessage("I'm Back :grinning:", channel.ID))
	// 	}

	case *slack.MessageEvent:

		//Identifying my ID (Bot Id)
		Name := strings.Fields(ev.Text)
		Word := Name[len(Name)-1]
		if cfg.Section("main").Key("name").String() == Word {
			log.Info("Msg: ", ev.Text)
		} else {
			log.Warning("Sent to User: ", Word)
			break Switching
		}

		//Regex to remove the @user
		reg, error := regexp.Compile(`\<.*\>`)
		if error != nil {
			log.Error("Failed to compile the message with regex")
		}
		message := reg.ReplaceAllString(ev.Text, "")
		if message == "" {
			log.Warn("Empty Message")
			break Switching
		}

		//Authentication section
		var ids bool
		chk := make(chan bool)
		if cfg.Section("main").Key("admin").String() == "true" {
			log.Info("Authorization Enabled !!!")
			ids = Auth(ev.User, cfg)
		} else if cfg.Section("main").Key("admin").String() == "false" {
			log.Info("Authorization Disabled !!!")
			ids = true
		} else {
			ids = true
		}

		if ids == true {
			//Notifying for long running jobs
			go func() {
				time.Sleep(7 * time.Second)
				select {
				case <-chk:
					runtime.Goexit()
				default:
					rtm.SendMessage(rtm.NewOutgoingMessage("`Wait..`:sleepy:", ev.Channel))
					runtime.Goexit()
				}
			}()
			// log.Info("Received:", ev.Text)
			rtm.SendMessage(rtm.NewOutgoingMessage(Format(message, cfg), ev.Channel))
			close(chk)
		} else if ids == false {
			rtm.SendMessage(rtm.NewOutgoingMessage("`Unauthorized` :cry:", ev.Channel))
		}

	case *slack.PresenceChangeEvent:
		log.Info("Presence Change :", ev)

	// case *slack.LatencyReport:
	//log.Info("Current latency :", ev.Value)

	case *slack.RTMError:
		log.Error("Error: ", ev.Error())

	case *slack.InvalidAuthEvent:
		log.Error("Invalid credentials")
		return

	default:
		log.Info("SlackBot: ", msg.Data)
	}
	wg.Done()
}
