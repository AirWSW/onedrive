package core

import "github.com/robfig/cron/v3"

func (od *OneDrive) Cron() error {
	c := cron.New(cron.WithSeconds())
	// Every half hour, starting a half hour from now
	c.AddFunc("@every 30m", func() {
		od.GetMicrosoftGraphAPIToken()
		od.SaveConfigFile()
	})
	c.AddFunc("@every 30m", func() {
		od.GetMicrosoftGraphAPIToken()
		od.SaveConfigFile()
	})
	c.Start()
	return nil
}
