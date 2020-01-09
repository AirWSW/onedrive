package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/AirWSW/onedrive/graphapi"
)

func (od *OneDrive) LoadDriveCacheFile() error {
	cacheFile := od.OneDriveDescription.DriveDescription.ID + ".cache.json"
	log.Println("Loading OneDrive cache file from " + cacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := ioutil.ReadFile(cacheFile)
	if _, ok := err.(*os.PathError); ok {
		log.Println("Creating OneDrive cache file " + cacheFile)
		return ioutil.WriteFile(cacheFile, []byte("{}"), 0644)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, od)
}

func (od *OneDrive) SaveDriveCacheFile() error {
	oneDriveCache := struct {
		DriveDescriptionCache        *graphapi.MicrosoftGraphDrive  `json:"driveDescriptionCache"`
		MicrosoftGraphDriveItemCache []MicrosoftGraphDriveItemCache `json:"microsoftGraphDriveItemCache"`
	}{
		od.OneDriveDescription.DriveDescription,
		od.MicrosoftGraphDriveItemCache,
	}

	cacheFile := od.OneDriveDescription.DriveDescription.ID + ".cache.json"
	bytes, err := json.Marshal(oneDriveCache)
	if err != nil {
		return err
	}

	log.Println("Saving OneDrive cache file to " + cacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	return ioutil.WriteFile(cacheFile, bytes, 0644)
}

func (od *OneDrive) DriveItemToCache(microsoftGraphDriveItem *graphapi.MicrosoftGraphDriveItem) (*MicrosoftGraphDriveItemCache, error) {
	parentReference := microsoftGraphDriveItem.ParentReference
	path := parentReference.Path + "/" + microsoftGraphDriveItem.Name
	cacheDescription := &CacheDescription{
		RequestURL:   path,
		Path:         path,
		LastUpdateAt: time.Now().Unix(),
		Status:       "Cached",
	}

	innerMicrosoftGraphDriveItemCache := []MicrosoftGraphDriveItemCache{}
	newMicrosoftGraphDriveItemCache := MicrosoftGraphDriveItemCache{}
	for _, children := range microsoftGraphDriveItem.Children {
		newMicrosoftGraphDriveItemCache = MicrosoftGraphDriveItemCache{
			CTag:                        children.CTag,
			Description:                 children.Description,
			File:                        children.File,
			Folder:                      children.Folder,
			Size:                        children.Size,
			ID:                          children.ID,
			CreatedAt:                   children.CreatedDateTime.Unix(),
			ETag:                        children.ETag,
			LastModifiedAt:              children.LastModifiedDateTime.Unix(),
			Name:                        children.Name,
			ParentReference:             children.ParentReference,
			WebURL:                      children.WebURL,
			AtMicrosoftGraphDownloadURL: children.AtMicrosoftGraphDownloadURL,
		}
		innerMicrosoftGraphDriveItemCache = append(innerMicrosoftGraphDriveItemCache, newMicrosoftGraphDriveItemCache)
	}

	microsoftGraphDriveItemCache := MicrosoftGraphDriveItemCache{
		CacheDescription:            cacheDescription,
		CTag:                        microsoftGraphDriveItem.CTag,
		Description:                 microsoftGraphDriveItem.Description,
		File:                        microsoftGraphDriveItem.File,
		Folder:                      microsoftGraphDriveItem.Folder,
		Size:                        microsoftGraphDriveItem.Size,
		Children:                    innerMicrosoftGraphDriveItemCache,
		ID:                          microsoftGraphDriveItem.ID,
		CreatedAt:                   microsoftGraphDriveItem.CreatedDateTime.Unix(),
		ETag:                        microsoftGraphDriveItem.ETag,
		LastModifiedAt:              microsoftGraphDriveItem.LastModifiedDateTime.Unix(),
		Name:                        microsoftGraphDriveItem.Name,
		ParentReference:             microsoftGraphDriveItem.ParentReference,
		WebURL:                      microsoftGraphDriveItem.WebURL,
		AtMicrosoftGraphDownloadURL: microsoftGraphDriveItem.AtMicrosoftGraphDownloadURL,
	}
	return &microsoftGraphDriveItemCache, nil
}

func (od *OneDrive) DriveItemCacheToPayLoad(microsoftGraphDriveItemCache *MicrosoftGraphDriveItemCache) (*DriveItemCachePayload, error) {
	innerDriveItemCachePayload := []DriveItemCachePayload{}
	newDriveItemCachePayload := DriveItemCachePayload{}
	for _, children := range microsoftGraphDriveItemCache.Children {
		newDriveItemCachePayload = DriveItemCachePayload{
			CTag:           children.CTag,
			Description:    children.Description,
			File:           children.File,
			Folder:         children.Folder,
			Size:           children.Size,
			ID:             children.ID,
			CreatedAt:      time.Unix(children.CreatedAt, 0),
			ETag:           children.ETag,
			LastModifiedAt: time.Unix(children.LastModifiedAt, 0),
			Name:           children.Name,
			DownloadURL:    nil,
		}
		innerDriveItemCachePayload = append(innerDriveItemCachePayload, newDriveItemCachePayload)
	}

	driveItemCachePayloadReference := &DriveItemCachePayloadReference{}
	parentReference := &graphapi.MicrosoftGraphItemReference{}
	relativePath := ""
	if microsoftGraphDriveItemCache.Children != nil {
		if microsoftGraphDriveItemCache.Children != nil || len(microsoftGraphDriveItemCache.Children) > 0 {
			parentReference = microsoftGraphDriveItemCache.Children[0].ParentReference
			relativePath = od.OneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
		} else {
			parentReference = microsoftGraphDriveItemCache.ParentReference
			relativePath = od.OneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
		}
	} else {
		parentReference = microsoftGraphDriveItemCache.ParentReference
		relativePath = od.OneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
	}
	driveItemCachePayloadReference = &DriveItemCachePayloadReference{
		DriveID:   parentReference.DriveID,
		DriveType: parentReference.DriveType,
		ID:        parentReference.ID,
		Path:      relativePath,
	}

	driveItemCachePayload := DriveItemCachePayload{
		LastUpdateAt:    time.Unix(microsoftGraphDriveItemCache.CacheDescription.LastUpdateAt, 0),
		CTag:            microsoftGraphDriveItemCache.CTag,
		Description:     microsoftGraphDriveItemCache.Description,
		File:            microsoftGraphDriveItemCache.File,
		Folder:          microsoftGraphDriveItemCache.Folder,
		Size:            microsoftGraphDriveItemCache.Size,
		Children:        innerDriveItemCachePayload,
		ID:              microsoftGraphDriveItemCache.ID,
		CreatedAt:       time.Unix(microsoftGraphDriveItemCache.CreatedAt, 0),
		ETag:            microsoftGraphDriveItemCache.ETag,
		LastModifiedAt:  time.Unix(microsoftGraphDriveItemCache.LastModifiedAt, 0),
		Name:            microsoftGraphDriveItemCache.Name,
		ParentReference: driveItemCachePayloadReference,
		DownloadURL:     nil,
	}
	return &driveItemCachePayload, nil
}

func (od *OneDrive) GetMicrosoftGraphDriveItemFromCache(str string) (*MicrosoftGraphDriveItemCache, error) {
	path, filename := RegularPathToPathFilename(str)
	str = od.OneDriveDescription.RelativePathToDriveRootPath(str)
	path = od.OneDriveDescription.RelativePathToDriveRootPath(path)
	log.Println("Hitting cache for " + path)
	for _, microsoftGraphDriveItemCache := range od.MicrosoftGraphDriveItemCache {
		cacheDescription := microsoftGraphDriveItemCache.CacheDescription
		if cacheDescription.Path == path {
			if time.Now().Unix()-cacheDescription.LastUpdateAt > 3600 {
				return nil, errors.New("MicrosoftGraphDriveItemCacheExpired")
			}
			if filename == "" {
				return &microsoftGraphDriveItemCache, nil
			} else {
				for _, children := range microsoftGraphDriveItemCache.Children {
					if children.Name == filename {
						if children.File != nil {
							children.CacheDescription = cacheDescription
							log.Println("HIT Cache "+cacheDescription.RequestURL, cacheDescription.LastUpdateAt)
							return &children, nil
						}
						for _, innerMicrosoftGraphDriveItemCache := range od.MicrosoftGraphDriveItemCache {
							innerCacheDescription := innerMicrosoftGraphDriveItemCache.CacheDescription
							if innerCacheDescription.Path == str {
								log.Println("HIT Cache "+innerCacheDescription.RequestURL, innerCacheDescription.LastUpdateAt)
								return &innerMicrosoftGraphDriveItemCache, nil
							}
						}
						log.Println("HIT Miss " + str)
						od.MicrosoftGraphDriveItemCache = append(od.MicrosoftGraphDriveItemCache, MicrosoftGraphDriveItemCache{
							CacheDescription: &CacheDescription{
								RequestURL:   str,
								Path:         str,
								LastUpdateAt: 0,
								Status:       "Wait",
							},
						})
						parentReference := microsoftGraphDriveItemCache.ParentReference
						newPath := parentReference.Path + "/" + microsoftGraphDriveItemCache.Name
						newChildren := children
						newChildren.ParentReference = &graphapi.MicrosoftGraphItemReference{
							DriveID:   parentReference.DriveID,
							DriveType: parentReference.DriveType,
							ID:        microsoftGraphDriveItemCache.ID,
							Path:      newPath,
						}
						cacheDescription := microsoftGraphDriveItemCache.CacheDescription
						newChildren.CacheDescription = &CacheDescription{
							RequestURL:   newPath,
							Path:         newPath,
							LastUpdateAt: cacheDescription.LastUpdateAt,
							Status:       "Wait",
						}
						return &newChildren, nil
					}
				}
				return nil, errors.New("NoMicrosoftGraphDriveItemCacheRecord")
			}
		}
	}
	od.MicrosoftGraphDriveItemCache = append(od.MicrosoftGraphDriveItemCache, MicrosoftGraphDriveItemCache{
		CacheDescription: &CacheDescription{
			RequestURL:   path,
			Path:         path,
			LastUpdateAt: 0,
			Status:       "Wait",
		},
	})
	return nil, errors.New("NoMicrosoftGraphDriveItemCacheRecord")
}

func (od *OneDrive) HitMicrosoftGraphDriveItemCache(path string) (*MicrosoftGraphDriveItemCache, error) {
	return od.GetMicrosoftGraphDriveItemFromCache(path)
}

func (od *OneDrive) GetMicrosoftGraphDriveItem(path string) (*DriveItemCachePayload, error) {
	microsoftGraphDriveItemCache, err := od.HitMicrosoftGraphDriveItemCache(path)
	go od.CronCacheMicrosoftGraphDrive()
	if err != nil {
		return nil, err
	}
	return od.DriveItemCacheToPayLoad(microsoftGraphDriveItemCache)
}

func (od *OneDrive) UpdateMicrosoftGraphDriveItemCache(cacheDescription *CacheDescription) (*MicrosoftGraphDriveItemCache, error) {
	return od.GetMicrosoftGraphAPIMeDriveChildrenRequest(cacheDescription.Path)
}
