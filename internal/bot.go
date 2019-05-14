package bot

import (
	"runtime"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
	"github.com/nlopes/slack"
)

//Run for running according to RTM events
func Run(msg slack.RTMEvent, wg *sync.WaitGroup, cfg *ini.File, channels []slack.Channel, rtm *slack.RTM) {
	switch ev := msg.Data.(type) {
	case *slack.HelloEvent:
		// Ignore hello

	case *slack.ConnectedEvent:
		// fmt.Println("Infos:", ev.Info)
		// fmt.Println("Connection counter:", ev.ConnectionCount)
		for _, channel := range channels {
			rtm.SendMessage(rtm.NewOutgoingMessage("I'm Back :grinning:", channel.ID))
		}

	case *slack.MessageEvent:
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
			go func() {
				time.Sleep(5 * time.Second)
				select {
				case <-chk:
					runtime.Goexit()
				default:
					rtm.SendMessage(rtm.NewOutgoingMessage("Wait..:sleepy:", ev.Channel))
					runtime.Goexit()
				}
			}()
			// log.Info("Received:", ev.Text)
			rtm.SendMessage(rtm.NewOutgoingMessage(Format(ev.Text, cfg), ev.Channel))
			close(chk)
		} else if ids == false {
			rtm.SendMessage(rtm.NewOutgoingMessage("Unauthorized :cry:", ev.Channel))
		}

	case *slack.PresenceChangeEvent:
		log.Info("Presence Change :", ev)

	case *slack.LatencyReport:
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
