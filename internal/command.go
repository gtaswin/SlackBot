package internal

import (
	"fmt"
	"io"
	"os/exec"
	"unicode/utf8"

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
		io.WriteString(stdin, text)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
	}
	a := string(out)
	// if a == "" {
	// 	a = "Done! "
	// }

	//Range management
	if utf8.RuneCountInString(a) >= 4000 {
		a = fmt.Sprintln(a[0:3500], "`...Out of range` :anguished:")
	}
	return a
}
