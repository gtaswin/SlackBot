package bot

import (
	"io"
	"log"
	"os/exec"

	"github.com/go-ini/ini"
)

//Command executing the script that present on ini file
func Command(text string, cfg *ini.File) string {
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
	if a == "" {
		a = "Done!"
	}
	return a
}
