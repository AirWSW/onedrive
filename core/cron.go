package core

import (
	"errors"
	"log"

	"github.com/AirWSW/onedrive/core/cache"
)

func (od *OneDrive) CronCacheMicrosoftGraphDrive() error {
	for i, microsoftGraphDriveItemCache := range od.DriveCacheCollection.MicrosoftGraphDriveItemCache {
		cacheDescription := microsoftGraphDriveItemCache.CacheDescription
		if err := cache.IsCacheNeedUpdate(&od.OneDriveDescription, cacheDescription); err != nil {
			log.Println("od.CronCacheMicrosoftGraphDrive", err)
			od.DriveCacheCollection.MicrosoftGraphDriveItemCache[i].CacheDescription.Status = "Caching"
			newMicrosoftGraphDriveItemCache, err := od.MicrosoftGraphAPI.UpdateMicrosoftGraphDriveItemCache(&od.OneDriveDescription, cacheDescription)
			if err != nil {
				log.Println("od.CronCacheMicrosoftGraphDrive", err)
				newMicrosoftGraphDriveItemCache = &microsoftGraphDriveItemCache
				newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Failed" // Failed, deleted
				od.DriveCacheCollection.MicrosoftGraphDriveItemCache[i] = *newMicrosoftGraphDriveItemCache
			} else {
				newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Cached"
				od.DriveCacheCollection.MicrosoftGraphDriveItemCache[i] = *newMicrosoftGraphDriveItemCache
				od.DriveCacheCollection.Save(od.OneDriveDescription.DriveDescription)
				return nil
			}
		}
	}
	return errors.New("od.CronCacheMicrosoftGraphDrive NothingNeedToCache")
}
