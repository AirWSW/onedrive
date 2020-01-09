package core

import (
	"log"
	"time"

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

func (od *OneDrive) CronCacheMicrosoftGraphDrive() error {
	mutex.Lock()
	defer mutex.Unlock()
	for i, microsoftGraphDriveItemCache := range od.MicrosoftGraphDriveItemCache {
		cacheDescription := microsoftGraphDriveItemCache.CacheDescription
		if time.Now().Unix()-cacheDescription.LastUpdateAt > od.OneDriveDescription.RefreshInterval && cacheDescription.Status != "Failed" {
			newMicrosoftGraphDriveItemCache, err := od.UpdateMicrosoftGraphDriveItemCache(cacheDescription)
			if err != nil {
				newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Failed"
				return err
			}
			newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Cached"
			od.MicrosoftGraphDriveItemCache[i] = *newMicrosoftGraphDriveItemCache
		}
	}
	return nil
}
