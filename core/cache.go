package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/core/utils"
	"github.com/AirWSW/onedrive/graphapi"
)

var mutex sync.Mutex

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
		LastUpdateAt: time.Unix(microsoftGraphDriveItemCache.CacheDescription.LastUpdateAt, 0).UTC(),
		DriveType:    parentReference.DriveType,
		Path:         relativePath,
	}

	innerDriveItemCachePayload := []DriveItemCachePayload{}
	newDriveItemCachePayload := DriveItemCachePayload{}
	for _, children := range microsoftGraphDriveItemCache.Children {
		innerDownloadURL := relativePath + "/" + children.Name
		innerDownloadURLPointer := &innerDownloadURL
		if children.Folder != nil {
			innerDownloadURLPointer = nil
		}
		newDriveItemCachePayload = DriveItemCachePayload{
			Description:    children.Description,
			File:           children.File,
			Folder:         children.Folder,
			Size:           children.Size,
			CreatedAt:      time.Unix(children.CreatedAt, 0).UTC(),
			LastModifiedAt: time.Unix(children.LastModifiedAt, 0).UTC(),
			Name:           children.Name,
			DownloadURL:    innerDownloadURLPointer,
		}
		innerDriveItemCachePayload = append(innerDriveItemCachePayload, newDriveItemCachePayload)
	}

	downloadURL := relativePath + "/" + microsoftGraphDriveItemCache.Name
	downloadURLPointer := &downloadURL
	if microsoftGraphDriveItemCache.Folder != nil {
		downloadURLPointer = nil
	}
	driveItemCachePayload := DriveItemCachePayload{
		Description:    microsoftGraphDriveItemCache.Description,
		File:           microsoftGraphDriveItemCache.File,
		Folder:         microsoftGraphDriveItemCache.Folder,
		Size:           microsoftGraphDriveItemCache.Size,
		Children:       innerDriveItemCachePayload,
		CreatedAt:      time.Unix(microsoftGraphDriveItemCache.CreatedAt, 0).UTC(),
		LastModifiedAt: time.Unix(microsoftGraphDriveItemCache.LastModifiedAt, 0).UTC(),
		Name:           microsoftGraphDriveItemCache.Name,
		Reference:      driveItemCachePayloadReference,
		DownloadURL:    downloadURLPointer,
	}
	return &driveItemCachePayload, nil
}

func (od *OneDrive) GetMicrosoftGraphDriveItemFromCache(str string) (*MicrosoftGraphDriveItemCache, error) {
	path, filename := utils.RegularPathToPathFilename(str)
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
						if children.Folder != nil {
							if children.Folder.ChildCount == 0 {
								children.CacheDescription = cacheDescription
								log.Println("HIT Cache "+cacheDescription.RequestURL, cacheDescription.LastUpdateAt)
								return &children, nil
							} else {
								for _, innerMicrosoftGraphDriveItemCache := range od.MicrosoftGraphDriveItemCache {
									innerCacheDescription := innerMicrosoftGraphDriveItemCache.CacheDescription
									if innerCacheDescription.Path == str {
										log.Println("HIT Cache "+innerCacheDescription.RequestURL, innerCacheDescription.LastUpdateAt)
										return &innerMicrosoftGraphDriveItemCache, nil
									}
								}
								log.Println("HIT Miss " + str)
								parentReference := microsoftGraphDriveItemCache.ParentReference
								newChildren := children
								newChildren.ParentReference = &graphapi.MicrosoftGraphItemReference{
									DriveID:   parentReference.DriveID,
									DriveType: parentReference.DriveType,
									ID:        microsoftGraphDriveItemCache.ID,
									Path:      str,
								}
								newChildren.CacheDescription = &CacheDescription{
									RequestURL:   str,
									Path:         str,
									LastUpdateAt: 0,
									Status:       "Wait",
								}
								od.MicrosoftGraphDriveItemCache = append(od.MicrosoftGraphDriveItemCache, newChildren)
								return &newChildren, nil
							}
						}
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
	newPath := utils.RegularPath(path)
	newPathLength := len(newPath)
	parentPath, filename := utils.RegularPathToPathFilename(path)
	driveVolumeMountRule := &description.DriveVolumeMount{}
	for _, driveVolumeMount := range od.OneDriveDescription.DriveVolumeMounts {
		target := utils.RegularPath(*driveVolumeMount.Target)
		targetLength := len(target)
		if newPathLength >= targetLength && newPath[0:targetLength] == target {
			newPath = utils.RegularPath(*driveVolumeMount.Source) + newPath[targetLength:newPathLength]
			driveVolumeMountRule = &driveVolumeMount
			continue
		}
	}

	microsoftGraphDriveItemCache, err := od.HitMicrosoftGraphDriveItemCache(newPath)
	go od.CronCacheMicrosoftGraphDrive()
	if err != nil {
		return nil, err
	}
	driveItemCachePayload, err := od.DriveItemCacheToPayLoad(microsoftGraphDriveItemCache)
	if err != nil {
		return nil, err
	}

	if newPath != utils.RegularPath(path) {
		if driveVolumeMountRule.Type != nil {
			if *driveVolumeMountRule.Type == "file.only" && driveItemCachePayload.File == nil {
				return nil, errors.New("file.only")
			}
		}
		driveItemCachePayload.Reference.Path = utils.RegularPath(parentPath)
		driveItemCachePayload.Name = filename
		innerDriveItemCachePayload := []DriveItemCachePayload{}
		newDriveItemCachePayload := DriveItemCachePayload{}
		for _, children := range driveItemCachePayload.Children {
			innerDownloadURL := od.OneDriveDescription.DriveRootPathToRelativePath(microsoftGraphDriveItemCache.ParentReference.Path)
			innerDownloadURL += "/" + children.Name
			innerDownloadURLPointer := &innerDownloadURL
			if children.Folder != nil {
				innerDownloadURLPointer = nil
			}
			newDriveItemCachePayload = DriveItemCachePayload{
				Description:    children.Description,
				File:           children.File,
				Folder:         children.Folder,
				Size:           children.Size,
				CreatedAt:      children.CreatedAt,
				LastModifiedAt: children.LastModifiedAt,
				Name:           children.Name,
				DownloadURL:    innerDownloadURLPointer,
			}
			innerDriveItemCachePayload = append(innerDriveItemCachePayload, newDriveItemCachePayload)
		}
	}

	return driveItemCachePayload, nil
}

func (od *OneDrive) DriveContentURLCacheToPayLoad(microsoftGraphDriveItemCache *MicrosoftGraphDriveItemCache) (*DriveItemCachePayload, error) {
	driveItemCachePayloadReference := &DriveItemCachePayloadReference{}
	parentReference := microsoftGraphDriveItemCache.ParentReference
	relativePath := od.OneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
	driveItemCachePayloadReference = &DriveItemCachePayloadReference{
		LastUpdateAt: time.Unix(microsoftGraphDriveItemCache.CacheDescription.LastUpdateAt, 0).UTC(),
		DriveType:    parentReference.DriveType,
		Path:         relativePath,
	}

	driveItemCachePayload := DriveItemCachePayload{
		Description:    microsoftGraphDriveItemCache.Description,
		File:           microsoftGraphDriveItemCache.File,
		Folder:         microsoftGraphDriveItemCache.Folder,
		Size:           microsoftGraphDriveItemCache.Size,
		CreatedAt:      time.Unix(microsoftGraphDriveItemCache.CreatedAt, 0).UTC(),
		LastModifiedAt: time.Unix(microsoftGraphDriveItemCache.LastModifiedAt, 0).UTC(),
		Name:           microsoftGraphDriveItemCache.Name,
		Reference:      driveItemCachePayloadReference,
		DownloadURL:    &microsoftGraphDriveItemCache.AtMicrosoftGraphDownloadURL,
	}
	return &driveItemCachePayload, nil
}

func (od *OneDrive) GetMicrosoftGraphDriveContentURLFromCache(str string) (*MicrosoftGraphDriveItemCache, error) {
	path, filename := utils.RegularPathToPathFilename(str)
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

func (od *OneDrive) HitMicrosoftGraphDriveContentURLCache(path string) (*MicrosoftGraphDriveItemCache, error) {
	return od.GetMicrosoftGraphDriveContentURLFromCache(path)
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveContentURL(path string) (*DriveItemCachePayload, error) {
	microsoftGraphDriveItemCache, err := od.HitMicrosoftGraphDriveContentURLCache(path)
	go od.CronCacheMicrosoftGraphDrive()
	if err != nil {
		return nil, err
	}
	return od.DriveContentURLCacheToPayLoad(microsoftGraphDriveItemCache)
}

func (od *OneDrive) UpdateMicrosoftGraphDriveItemCache(cacheDescription *CacheDescription) (*MicrosoftGraphDriveItemCache, error) {
	return od.GetMicrosoftGraphAPIMeDriveChildrenRequest(cacheDescription.Path)
}
