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
			rtm.SendMessage(rtm.NewOutgoingMessage("I'm Back", channel.ID))
		}

	case *slack.MessageEvent:
		ids := Auth(ev.User, cfg)
		chk := make(chan bool)

		if ids == true {
			go func() {
				time.Sleep(5 * time.Second)
				select {
				case <-chk:
					runtime.Goexit()
				default:
					rtm.SendMessage(rtm.NewOutgoingMessage("Still Running...", ev.Channel))
					runtime.Goexit()
				}
			}()

			rtm.SendMessage(rtm.NewOutgoingMessage(Format(ev.Text, cfg), ev.Channel))
			close(chk)
		}

	case *slack.PresenceChangeEvent:
		log.Info("Presence Change :", ev)

	case *slack.LatencyReport:
		log.Info("Current latency :", ev.Value)

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
