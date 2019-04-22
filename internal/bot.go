package bot

import (
  "fmt"
  "sync"
  "github.com/go-ini/ini"
  "github.com/nlopes/slack"
)

func Run(msg slack.RTMEvent, wg *sync.WaitGroup, cfg *ini.File, channels []slack.Channel, rtm *slack.RTM){
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
      ids := auth(ev.User, cfg)
      if ids == true {
      rtm.SendMessage(rtm.NewOutgoingMessage(format(ev.Text, cfg), ev.Channel))
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
