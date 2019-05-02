package main

import (
  "fmt"
  "flag"
  "strconv"
  "sync"
  "github.com/gtaswin/Slackbot/internal"
  "github.com/go-ini/ini"
  "github.com/nlopes/slack"
)

var (
	conf = flag.String("conf", "config/config.ini", "Configuration File")
)

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
  go bot.Run(msg, &wg, cfg, channels, rtm)
}
}
