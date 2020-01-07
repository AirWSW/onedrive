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
	// Every minuite, starting a minuite from now
	c.AddFunc("@every 1m", func() {
		od.CronCacheDrive()
		od.SaveConfigFile()
	})
	// Every minuite, starting a minuite from now
	c.AddFunc("@every 1m", func() {
		od.CronCacheDrivePathContentURL()
		od.SaveConfigFile()
	})
	c.Start()
	return nil
}

func (od *OneDrive) CronCacheDrive() error {
	for i, item := range od.DriveCacheContentURL {
		if time.Now().Sub(item.UpdateAt) < time.Second*1200 {
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
		if time.Now().Sub(item.UpdateAt) < time.Second*600 {
			driveCache, err := od.getDriveCache(item.Path, item.RequestURL)
			if err != nil {
				return err
			}
			od.DriveCache[i] = *driveCache
		}
	}
	return nil
}
