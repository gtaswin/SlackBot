package bot

import (
	"fmt"
	"runtime"
	"sync"
	"time"

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
		fmt.Printf("Presence Change: %v\n", ev)

	case *slack.LatencyReport:
		fmt.Printf("Current latency: %v\n", ev.Value)

	case *slack.RTMError:
		fmt.Printf("Error: %s\n", ev.Error())

	case *slack.InvalidAuthEvent:
		fmt.Printf("Invalid credentials")
		return

	default:
		fmt.Printf("SlackBot: %v\n", msg.Data)
	}
	wg.Done()
}
