package bot

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	log "github.com/Sirupsen/logrus"
	"github.com/go-ini/ini"
)

//Format for parsing the given text
func Format(text string, cfg *ini.File) string {
	var out string
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		log.Error(err)
	}

	text = reg.ReplaceAllString(text, "")
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	all := cfg.Section("chat").KeyStrings()

	if cfg.Section("chat").Key(text).String() == "" {
		for _, check := range all {
			s := Suggestion(text, check)
			if s == true {
				sug := fmt.Sprintln("Do you mean: ", check)
				return sug
			}
		}
	}
	if cfg.Section("chat").Key(text).String() == "" {
		a := cfg.Section("main").Key("noentry").String()
		return a
	}
	array := strings.Fields(text)
	if array[0] == cfg.Section("main").Key("command").String() {
		out = Command(text, cfg)
	} else {
		out = cfg.Section("chat").Key(text).String()
	}
	if utf8.RuneCountInString(out) >= 3000 {
		out = "Out of range"
	}
	return out
}
