package core

import (
	"time"
)

func (od *OneDrive) CronCacheMicrosoftGraphDrive() error {
	mutex.Lock()
	defer mutex.Unlock()
	for i, microsoftGraphDriveItemCache := range od.MicrosoftGraphDriveItemCache {
		cacheDescription := microsoftGraphDriveItemCache.CacheDescription
		if time.Now().Unix()-cacheDescription.LastUpdateAt > od.OneDriveDescription.RefreshInterval && cacheDescription.Status != "Failed" {
			newMicrosoftGraphDriveItemCache, err := od.UpdateMicrosoftGraphDriveItemCache(cacheDescription)
			if err != nil {
				newMicrosoftGraphDriveItemCache = &microsoftGraphDriveItemCache
				newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Failed"
				od.MicrosoftGraphDriveItemCache[i] = *newMicrosoftGraphDriveItemCache
				return err
			}
			newMicrosoftGraphDriveItemCache.CacheDescription.Status = "Cached"
			od.MicrosoftGraphDriveItemCache[i] = *newMicrosoftGraphDriveItemCache
		}
	}
	return nil
}
