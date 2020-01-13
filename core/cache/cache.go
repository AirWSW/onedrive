package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/AirWSW/onedrive/core/utils"
	"github.com/AirWSW/onedrive/graphapi"
)

var mutex sync.Mutex

func DriveItemToCache(microsoftGraphDriveItem *graphapi.MicrosoftGraphDriveItem) (*MicrosoftGraphDriveItemCache, error) {
	parentReference := microsoftGraphDriveItem.ParentReference
	path := parentReference.Path + "/" + microsoftGraphDriveItem.Name
	if parentReference.Path == "" {
		path = "/drive/root:"
	}
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

func IsCacheInvalid(odd oneDriveDescription, cacheDescription *CacheDescription) error {
	// log.Println("IsCacheInvalid " + cacheDescription.Status + " " + cacheDescription.Path + " will expired at " + time.Unix(cacheDescription.LastUpdateAt, 0).UTC().String())
	if cacheDescription.Status == "Wait" {
		return errors.New("MicrosoftGraphDriveItemCacheStatusWait " + cacheDescription.Path)
	}
	if cacheDescription.Status == "Caching" {
		return errors.New("MicrosoftGraphDriveItemCacheStatusCaching " + cacheDescription.Path)
	}
	if cacheDescription.Status == "Failed" {
		return errors.New("MicrosoftGraphDriveItemCacheStatusFailed " + cacheDescription.Path)
	}
	if time.Now().Unix()-cacheDescription.LastUpdateAt > graphapi.AtMicrosoftGraphDownloadURLAvailablePeriod {
		return errors.New("MicrosoftGraphDriveItemCacheExpired " + cacheDescription.Path + " " + time.Unix(cacheDescription.LastUpdateAt, 0).UTC().String())
	}
	return nil
}

func IsCacheNeedUpdate(odd oneDriveDescription, cacheDescription *CacheDescription) error {
	// log.Println("IsCacheNeedUpdate " + cacheDescription.Status + " " + cacheDescription.Path + " will expired at " + time.Unix(cacheDescription.LastUpdateAt, 0).UTC().String())
	if cacheDescription.Status == "Wait" {
		return errors.New("MicrosoftGraphDriveItemCacheStatusWait " + cacheDescription.Path)
	}
	if cacheDescription.Status == "Caching" {
		return nil
	}
	if cacheDescription.Status == "Failed" {
		return nil
	}
	if time.Now().Unix()-cacheDescription.LastUpdateAt > graphapi.AtMicrosoftGraphDownloadURLAvailableSafePeriod-odd.GetRefreshInterval() {
		return errors.New("MicrosoftGraphDriveItemCacheNeedUpdate " + cacheDescription.Path + " " + time.Unix(cacheDescription.LastUpdateAt, 0).UTC().String())
	}
	return nil
}

type oneDriveDescription interface {
	GetRefreshInterval() int64
	RelativePathToDriveRootPath(string) string
}

func (dcc *DriveCacheCollection) Load(microsoftGraphDrive *graphapi.MicrosoftGraphDrive) error {
	cacheFile := microsoftGraphDrive.ID + ".cache.json"
	log.Println("Loading OneDrive cache file from", cacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := ioutil.ReadFile(cacheFile)
	if _, ok := err.(*os.PathError); ok {
		log.Println("Creating OneDrive cache file", cacheFile)
		return ioutil.WriteFile(cacheFile, []byte("{}"), 0644)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dcc)
}

func (dcc *DriveCacheCollection) Save(microsoftGraphDrive *graphapi.MicrosoftGraphDrive) error {
	oneDriveCache := struct {
		DriveDescriptionCache        *graphapi.MicrosoftGraphDrive  `json:"driveDescriptionCache"`
		MicrosoftGraphDriveItemCache []MicrosoftGraphDriveItemCache `json:"microsoftGraphDriveItemCache"`
	}{
		microsoftGraphDrive,
		dcc.MicrosoftGraphDriveItemCache,
	}

	cacheFile := microsoftGraphDrive.ID + ".cache.json"
	bytes, err := json.Marshal(oneDriveCache)
	if err != nil {
		return err
	}

	log.Println("Saving OneDrive cache file to", cacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	return ioutil.WriteFile(cacheFile, bytes, 0644)
}

func (dcc *DriveCacheCollection) HitMicrosoftGraphDriveItemCache(odd oneDriveDescription, path string) (*MicrosoftGraphDriveItemCache, error) {
	return dcc.GetMicrosoftGraphDriveItemFromCache(odd, path)
}

func (dcc *DriveCacheCollection) GetMicrosoftGraphDriveItemFromCache(odd oneDriveDescription, path string) (*MicrosoftGraphDriveItemCache, error) {
	return dcc.GetMicrosoftGraphDriveFromCache(odd, path, false)
}

func (dcc *DriveCacheCollection) HitMicrosoftGraphDriveContentURLCache(odd oneDriveDescription, path string) (*MicrosoftGraphDriveItemCache, error) {
	return dcc.GetMicrosoftGraphDriveContentURLFromCache(odd, path)
}

func (dcc *DriveCacheCollection) GetMicrosoftGraphDriveContentURLFromCache(odd oneDriveDescription, path string) (*MicrosoftGraphDriveItemCache, error) {
	return dcc.GetMicrosoftGraphDriveFromCache(odd, path, true)
}

func (dcc *DriveCacheCollection) GetMicrosoftGraphDriveFromCache(odd oneDriveDescription, path string, isContentURL bool) (*MicrosoftGraphDriveItemCache, error) {
	subPath, filename := utils.RegularPathToPathFilename(path)
	path = odd.RelativePathToDriveRootPath(path)
	subPath = odd.RelativePathToDriveRootPath(subPath)
	if path == "/drive/root:" {
		subPath, filename = "/drive/root:", ""
	}
	log.Println("Hitting cache for", subPath)
	for _, microsoftGraphDriveItemCache := range dcc.MicrosoftGraphDriveItemCache {
		cacheDescription := microsoftGraphDriveItemCache.CacheDescription
		if cacheDescription.Path == subPath {
			if err := IsCacheInvalid(odd, cacheDescription); err != nil {
				// dcc.MicrosoftGraphDriveItemCache[i].CacheDescription.Status = "Failed"
				return nil, err
			}
			if filename == "" {
				log.Println("Cache hitted for", cacheDescription.RequestURL, time.Unix(cacheDescription.LastUpdateAt, 0).UTC())
				return &microsoftGraphDriveItemCache, nil
			} else {
				return dcc.GetMicrosoftGraphDriveFromCacheStep2(odd, microsoftGraphDriveItemCache, filename, path, isContentURL)
			}
		}
	}
	dcc.MicrosoftGraphDriveItemCache = append(dcc.MicrosoftGraphDriveItemCache, MicrosoftGraphDriveItemCache{
		CacheDescription: &CacheDescription{
			RequestURL:   subPath,
			Path:         subPath,
			LastUpdateAt: 0,
			Status:       "Wait",
		},
	})
	return nil, errors.New("NoMicrosoftGraphDriveItemCacheRecord " + subPath)
}

func (dcc *DriveCacheCollection) GetMicrosoftGraphDriveFromCacheStep2(odd oneDriveDescription, this MicrosoftGraphDriveItemCache, filename, path string, isContentURL bool) (*MicrosoftGraphDriveItemCache, error) {
	cacheDescription := this.CacheDescription
	for _, children := range this.Children {
		if children.Name == filename {
			if children.File != nil {
				children.CacheDescription = cacheDescription
				log.Println("Cache hitted for", cacheDescription.RequestURL, time.Unix(cacheDescription.LastUpdateAt, 0).UTC())
				return &children, nil
			}
			if !isContentURL {
				return dcc.GetMicrosoftGraphDriveFromCacheStep3(odd, this, children, path)
			}
		}
	}
	return nil, errors.New("NoMicrosoftGraphDriveItemCacheRecord " + path)
}

func (dcc *DriveCacheCollection) GetMicrosoftGraphDriveFromCacheStep3(odd oneDriveDescription, this MicrosoftGraphDriveItemCache, children MicrosoftGraphDriveItemCache, path string) (*MicrosoftGraphDriveItemCache, error) {
	cacheDescription := this.CacheDescription
	if children.Folder != nil {
		if children.Folder.ChildCount == 0 {
			children.CacheDescription = cacheDescription
			log.Println("Cache hitted for", cacheDescription.RequestURL, time.Unix(cacheDescription.LastUpdateAt, 0).UTC())
			return &children, nil
		} else {
			var err error = nil
			for _, innerMicrosoftGraphDriveItemCache := range dcc.MicrosoftGraphDriveItemCache {
				innerCacheDescription := innerMicrosoftGraphDriveItemCache.CacheDescription
				if innerCacheDescription.Path == path {
					if err = IsCacheInvalid(odd, innerCacheDescription); err != nil {
						// return &this, err
					} else {
						log.Println("Cache hitted for", innerCacheDescription.RequestURL, time.Unix(innerCacheDescription.LastUpdateAt, 0).UTC())
						return &innerMicrosoftGraphDriveItemCache, nil
					}
				}
			}
			log.Println("Cache missed for", path)
			newChildren := children
			newChildren.CacheDescription = &CacheDescription{
				RequestURL:   path,
				Path:         path,
				LastUpdateAt: 0,
				Status:       "Wait",
			}
			if err == nil {
				dcc.MicrosoftGraphDriveItemCache = append(dcc.MicrosoftGraphDriveItemCache, newChildren)
			}
			return &newChildren, err
		}
	}
	return nil, errors.New("NoMicrosoftGraphDriveItemCacheRecord " + path)
}
