package core

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func (od *OneDrive) Cron() error {
	c := cron.New(cron.WithSeconds())
	// Every half hour, starting a half hour from now
	c.AddFunc("@every 30m", func() {
		od.GetMicrosoftGraphAPIToken()
		od.SaveConfigFile()
	})
	// Every 10 minuites, starting 10 minuites from now
	c.AddFunc("@every 60s", func() {
		log.Println("Start CronCacheDriveItems")
		od.CronCacheDriveItems()
		od.SaveDriveCacheFile()
	})
	c.Start()
	return nil
}

func (od *OneDrive) CronCacheDriveItems() error {
	for i, driveItemsCache := range od.DriveItemsCaches {
		driveItemsReference := driveItemsCache.DriveItemsReference
		if time.Now().Unix()-driveItemsReference.LastUpdateAt > 1800 {
			reqURL := od.DriveDescriptionConfig.EndPointURI
			reqURL += "/me"
			reqURL += driveItemsReference.Path
			reqURL += ":/children"
			log.Println("Updating: " + reqURL)
			graphAPIDriveItems, err := od.GetGraphAPIDriveItemsRequest(reqURL)
			if err != nil {
				return err
			}
			driveItemsCache, err := od.GraphAPIDriveItemsToDriveItemsCache(graphAPIDriveItems)
			if err != nil {
				return err
			}
			od.DriveItemsCaches[i] = *driveItemsCache
		}
	}
	return nil
}
