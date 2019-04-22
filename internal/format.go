package bot

import (
  "strings"
  "log"
  "regexp"
  "github.com/go-ini/ini"
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
