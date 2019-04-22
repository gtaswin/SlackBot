package bot

import (
  "fmt"
  "github.com/go-ini/ini"
)


func auth(id string, cfg *ini.File) bool {
  var a bool

  user := cfg.Section("admin").Keys()

  for _, users := range user {
  c := fmt.Sprintf("%v", users)
  if c == id {
  a = true
  }
  }

  return a
}
