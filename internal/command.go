package Internal

import (
	"fmt"
	"io"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

//Command executing the script that present on ini file
func Command(text string, cfg *ini.File) string {
	cmd := exec.Command(cfg.Section("main").Key("shell").String())
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Error(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, cfg.Section("chat").Key(text).String())
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
	}
	a := string(out)
	if a == "" {
		a = "Done! "
	} else if a != "" {
		a = fmt.Sprintln("Here the Output :sunglasses:\n", a)
	}
	return a
}
