package main

import (
  "fmt"
  "strings"
  "log"
  "io"
  "flag"
  "strconv"
  "sync"
  "os/exec"
  "regexp"
  "github.com/gtaswin/Slackbot/internal"
  "github.com/go-ini/ini"
  "github.com/nlopes/slack"
)

var (
	conf = flag.String("conf", "config/config.ini", "Configuration File")
)

func format(text string, cfg *ini.File) string {
  var out string
  reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
      if err != nil {
          log.Print(err)
      }

  text = reg.ReplaceAllString(text, "")
  text = strings.TrimSpace(text)
  text = strings.ToLower(text)
  array := strings.Fields(text)
  if array[0] == cfg.Section("main").Key("command").String() {
  out = command(text, cfg)
  } else {
  out = cfg.Section("chat").Key(text).String()
  }

  return out
}

func command(text string, cfg *ini.File) string {
	cmd := exec.Command(cfg.Section("main").Key("shell").String())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Print(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, cfg.Section("chat").Key(text).String())
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	}
	a := string(out)
  return a
}

func bot(msg slack.RTMEvent, wg *sync.WaitGroup, cfg *ini.File, channels []slack.Channel, rtm *slack.RTM){
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

func main() {
	flag.Parse()

	cfg, err := ini.Load(*conf)
	if err != nil {
        panic("Fail to read Configuration")
    }
	token := cfg.Section("main").Key("token").String()
  fmt.Println(token)
	debug, _ := strconv.ParseBool(cfg.Section("main").Key("debug").String())
	api := slack.New(token, slack.OptionDebug(debug),)
	// api = slack.New(token, slack.OptionDebug(false), slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),)
	rtm := api.NewRTM()
	channels, _ := api.GetChannels(false)

  go rtm.ManageConnection()

  var wg sync.WaitGroup
  for msg := range rtm.IncomingEvents {
  cfg, err := ini.Load(*conf)
  if err != nil {
        panic("Fail to read Configuration")
  }
  wg.Add(1)
  go bot(msg, &wg, cfg, channels, rtm)
}
}
