package internal

import (
	"github.com/go-ini/ini"
	"github.com/robfig/cron"
)

// Cron function to schedule the job and sent the output as the message.
func Cron(cfg *ini.File, Cronchan chan string) {
	c := cron.New()
	c.Start()
	for _, a := range cfg.Section("cron").KeyStrings() {
		b := a
		c.AddFunc(a, func() { Cronchan <- Command(cfg.Section("cron").Key(b).String(), cfg) })
	}
}
