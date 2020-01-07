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
	c.AddFunc("@every 10s", func() {
		log.Println("Start CronCacheDrive")
		od.CronCacheDrive()
		od.SaveDriveCacheFile()
	})
	// Every 20 minuites, starting 20 minuites from now
	c.AddFunc("@every 10s", func() {
		log.Println("Start CronCacheDrivePathContentURL")
		od.CronCacheDrivePathContentURL()
		od.SaveDriveCacheFile()
	})
	c.Start()
	return nil
}

func (od *OneDrive) CronCacheDrive() error {
	for i, item := range od.DriveCacheContentURL {
		if item.UpdateAt-time.Now().Unix() > 1200 {
			log.Println("Updating: " + item.RequestURL)
			driveCacheContentURL, err := od.getDriveCacheContentURL(item.RequestURL)
			if err != nil {
				return err
			}
			od.DriveCacheContentURL[i] = *driveCacheContentURL
		}
	}
	return nil
}

func (od *OneDrive) CronCacheDrivePathContentURL() error {
	for i, item := range od.DriveCache {
		if item.UpdateAt-time.Now().Unix() > 600 {
			log.Println("Updating: " + item.RequestURL)
			driveCache, err := od.getDriveCache(item.Path, item.RequestURL)
			if err != nil {
				return err
			}
			od.DriveCache[i] = *driveCache
		}
	}
	return nil
}
