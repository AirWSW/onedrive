package collection

import (
	"fmt"
	"log"

	"github.com/robfig/cron/v3"
)

func (odc *OneDriveCollection) CronStartAll() error {
	c := cron.New(cron.WithSeconds())
	// Every half hour, starting a half hour from now
	log.Printf("@every 30m api.RefreshMicrosoftGraphAPIToken\n")
	c.AddFunc("@every 30m", func() {
		for _, oneDrive := range odc.OneDrives {
			oneDrive.MicrosoftGraphAPI.RefreshMicrosoftGraphAPIToken()
		}
		odc.SaveConfigFile()
	})
	for _, oneDrive := range odc.OneDrives {
		refreshInterval := oneDrive.OneDriveDescription.GetRefreshInterval()
		log.Printf("@every %ds od.CronCacheMicrosoftGraphDrive\n", refreshInterval)
		c.AddFunc(fmt.Sprintf("@every %ds", refreshInterval), func() {
			// log.Printf("start @every %ds od.CronCacheMicrosoftGraphDrive\n", refreshInterval)
			// defer log.Printf("end @every %ds od.CronCacheMicrosoftGraphDrive\n", refreshInterval)
			if err := oneDrive.CronCacheMicrosoftGraphDrive(); err != nil {
				// log.Printf("@every %ds od.CronCacheMicrosoftGraphDrive %v\n", refreshInterval, err)
			} else {
				// oneDrive.DriveCacheCollection.Save(oneDrive.OneDriveDescription.DriveDescription)
			}
		})
	}
	c.Start()
	return nil
}
