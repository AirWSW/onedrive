package collection

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

func (odc *OneDriveCollection) CronStartAll() error {
	c := cron.New(cron.WithSeconds())
	// Every half hour, starting a half hour from now
	c.AddFunc("@every 30m", func() {
		for _, oneDrive := range odc.OneDrives {
			oneDrive.MicrosoftGraphAPI.RefreshMicrosoftGraphAPIToken()
		}
		odc.SaveConfigFile()
	})
	for _, oneDrive := range odc.OneDrives {
		refreshInterval := oneDrive.OneDriveDescription.GetRefreshInterval()
		c.AddFunc(fmt.Sprintf("@every %ds", refreshInterval), func() {
			log.Printf("start @every %ds odc.CronCacheMicrosoftGraphDriveItem\n", refreshInterval)
			defer log.Printf("end @every %ds odc.CronCacheMicrosoftGraphDriveItem\n", refreshInterval)
			if err := oneDrive.CronCacheMicrosoftGraphDrive(); err == nil {
				oneDrive.DriveCacheCollection.Save(oneDrive.OneDriveDescription.DriveDescription)
			} else {
				log.Printf("@every %ds odc.CronCacheMicrosoftGraphDriveItem %v\n", refreshInterval, err)
			}
		})
	}
	c.Start()
	return nil
}
