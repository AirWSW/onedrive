package core

import (
	"errors"
	"log"
	"time"

	"github.com/AirWSW/onedrive/core/cache"
	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/core/utils"
	"github.com/AirWSW/onedrive/graphapi"
)

func (od *OneDrive) GetMicrosoftGraphDriveItem(path string) (*DriveItemCachePayload, error) {
	newPath := utils.RegularPath(path)
	newPathLength := len(newPath)
	parentPath, filename := utils.RegularPathToPathFilename(path)
	if newPath == "/drive/root:" {
		newPath = "/drive/root"
		parentPath, filename = "/drive/root", ""
	}
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

	microsoftGraphDriveItemCache, err := od.DriveCacheCollection.HitMicrosoftGraphDriveItemCache(&od.OneDriveDescription, newPath)
	if err != nil {
		go func() {
			if err := od.CronCacheMicrosoftGraphDrive(); err != nil {
				log.Println(err)
			}
		}()
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

func (od *OneDrive) DriveItemCacheToPayLoad(microsoftGraphDriveItemCache *cache.MicrosoftGraphDriveItemCache) (*DriveItemCachePayload, error) {
	oneDriveDescription := od.OneDriveDescription
	driveItemCachePayloadReference := &DriveItemCachePayloadReference{}
	parentReference := &graphapi.MicrosoftGraphItemReference{}
	relativePath := ""
	if microsoftGraphDriveItemCache.Children != nil && len(microsoftGraphDriveItemCache.Children) > 0 {
		// Folder item which has children items
		parentReference = microsoftGraphDriveItemCache.Children[0].ParentReference
		relativePath = oneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
	} else if microsoftGraphDriveItemCache.Folder != nil && microsoftGraphDriveItemCache.Children == nil {
		// Folder item which does NOT have children items
		parentReference = microsoftGraphDriveItemCache.ParentReference
		relativePath = oneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
	} else if microsoftGraphDriveItemCache.File != nil {
		// File item
		parentReference = microsoftGraphDriveItemCache.ParentReference
		relativePath = oneDriveDescription.DriveRootPathToRelativePath(parentReference.Path)
	} else {
		// Folder item which does NOT sync children items yet
		parentReference = microsoftGraphDriveItemCache.ParentReference
		relativePath = oneDriveDescription.DriveRootPathToRelativePath(parentReference.Path + "/" + microsoftGraphDriveItemCache.Name)
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

	name := "root"
	if driveItemCachePayloadReference.Path != "/" {
		name = microsoftGraphDriveItemCache.Name
	} else if oneDriveDescription.OneDriveName != nil {
		name = *oneDriveDescription.OneDriveName
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
		Name:           name,
		Reference:      driveItemCachePayloadReference,
		DownloadURL:    downloadURLPointer,
	}
	return &driveItemCachePayload, nil
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveContentURL(path string) (*DriveItemCachePayload, error) {
	microsoftGraphDriveItemCache, err := od.DriveCacheCollection.HitMicrosoftGraphDriveContentURLCache(&od.OneDriveDescription, path)
	if err != nil {
		go func() {
			if err := od.CronCacheMicrosoftGraphDrive(); err != nil {
				log.Println(err)
			}
		}()
		return nil, err
	}
	return od.DriveContentURLCacheToPayLoad(microsoftGraphDriveItemCache)
}

func (od *OneDrive) DriveContentURLCacheToPayLoad(microsoftGraphDriveItemCache *cache.MicrosoftGraphDriveItemCache) (*DriveItemCachePayload, error) {
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
		DownloadURL:    microsoftGraphDriveItemCache.AtMicrosoftGraphDownloadURL,
	}
	return &driveItemCachePayload, nil
}