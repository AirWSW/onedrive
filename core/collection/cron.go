package collection

import (
	"log"

	"github.com/robfig/cron/v3"
)

func (odc *OneDriveCollection) CronStart() error {
	c := cron.New(cron.WithSeconds())
	// Every half hour, starting a half hour from now
	c.AddFunc("@every 30m", func() {
		for _, oneDrive := range odc.OneDrives {
			oneDrive.MicrosoftGraphAPI.RefreshMicrosoftGraphAPIToken()
		}
		odc.SaveConfigFile()
	})
	// Every minuite, starting a minuite from now
	c.AddFunc("@every 60s", func() {
		log.Println("Start CronCacheMicrosoftGraphDriveItem")
		for _, oneDrive := range odc.OneDrives {
			oneDrive.CronCacheMicrosoftGraphDrive()
			oneDrive.SaveDriveCacheFile()
		}
	})
	c.Start()
	return nil
}
