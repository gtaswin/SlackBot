package bot

import (
  "log"
  "io"
  "os/exec"
  "github.com/go-ini/ini"
)

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
