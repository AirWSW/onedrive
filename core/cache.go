package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GraphAPIDriveItems struct {
	OdataNextLink *string             `json:"@odata.nextLink"`
	DriveItems    []GraphAPIDriveItem `json:"value"`
}

type GraphAPIDriveItem struct {
	ID                        string                           `json:"id"`
	MicrosoftGraphDownloadURL *string                          `json:"@microsoft.graph.downloadUrl"`
	Name                      string                           `json:"name"`
	Size                      int                              `json:"size"`
	WebURL                    string                           `json:"webUrl"`
	CreatedDateTime           string                           `json:"createdDateTime"`
	LastModifiedDateTime      string                           `json:"lastModifiedDateTime"`
	ParentReference           GraphAPIDriveItemParentReference `json:"parentReference"`
	Folder                    *GraphAPIDriveItemFolder         `json:"folder"`
	File                      *GraphAPIDriveItemFile           `json:"file"`
}

type GraphAPIDriveItemParentReference struct {
	ID        string `json:"id"`
	DriveID   string `json:"driveId"`
	DriveType string `json:"driveType"`
	Path      string `json:"path"`
}

type GraphAPIDriveItemFolder struct {
	ChildCount int `json:"childCount"`
}

type GraphAPIDriveItemFile struct {
	MimeType string                      `json:"mimeType"`
	Hashes   GraphAPIDriveItemFileHashes `json:"hashes"`
}

type GraphAPIDriveItemFileHashes struct {
	QuickXorHash string  `json:"quickXorHash"`
	Sha1Hash     *string `json:"sha1Hash"`
}

func (od *OneDrive) PathToGraphAPIDriveItemsRequestURL(path string) string {
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += "/me"
	reqURL += RegularRootPath(od.DriveDescriptionConfig.RootPath)
	reqURL += RegularPath(path)
	reqURL += ":/children"
	return reqURL
}

func (od *OneDrive) PathToGraphAPIDriveItemsPath(path string) string {
	newPath := RegularRootPath(od.DriveDescriptionConfig.RootPath)
	newPath += RegularPath(path)
	return newPath
}

func (od *OneDrive) GraphAPIDriveItemsPathToPath(gPath string) string {
	rootPath := RegularRootPath(od.DriveDescriptionConfig.RootPath)
	gPath, _ = url.QueryUnescape(gPath)
	strS := strings.Split(gPath, rootPath)
	path := strS[0]
	for i, str := range strS {
		if i == 1 {
			path += str
		} else if i > 1 {
			path += rootPath + str
		}
	}
	if path == "" {
		path += "/"
	}
	return path
}

func (od *OneDrive) GetGraphAPIDriveItems(path string) (*GraphAPIDriveItems, error) {
	reqURL := od.PathToGraphAPIDriveItemsRequestURL(path)
	graphAPIDriveItems, err := od.GetGraphAPIDriveItemsRequest(reqURL)
	if err != nil {
		return nil, err
	}

	log.Printf("od.GetGraphAPIDriveItems get %d items for %s\n", len(graphAPIDriveItems.DriveItems), reqURL)
	return graphAPIDriveItems, nil
}

func (od *OneDrive) GetGraphAPIDriveItemsRequest(reqURL string) (*GraphAPIDriveItems, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+od.MicrosoftGraphAPIToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	graphAPIDriveItems := &GraphAPIDriveItems{}
	if err = json.Unmarshal([]byte(body), graphAPIDriveItems); err != nil {
		return nil, err
	}
	if graphAPIDriveItems.DriveItems == nil || len(graphAPIDriveItems.DriveItems) == 0 {
		log.Println(string(body))
		// return nil, errors.New("GraphAPIRequestError")
	}

	log.Println("od.GetGraphAPIDriveItemsRequest " + reqURL)

	if graphAPIDriveItems.OdataNextLink != nil {
		newGraphAPIDriveItems, err := od.GetGraphAPIDriveItemsRequest(*graphAPIDriveItems.OdataNextLink)
		if err != nil {
			return graphAPIDriveItems, err
		}
		for _, driveItem := range newGraphAPIDriveItems.DriveItems {
			graphAPIDriveItems.DriveItems = append(graphAPIDriveItems.DriveItems, driveItem)
		}
	}

	return graphAPIDriveItems, nil
}

// DriveItemsCache presents drive items cache structure
type DriveItemsCache struct {
	DriveItems          []DriveItemCache         `json:"driveItem"`
	DriveItemsReference DriveItemsReferenceCache `json:"driveItemReference"`
}

type DriveItemCache struct {
	ID             string  `json:"id"`
	DownloadURL    *string `json:"downloadUrl"`
	Name           string  `json:"name"`
	Size           int     `json:"size"`
	ChildCount     *int    `json:"childCount"`
	MimeType       *string `json:"mimeType"`
	QuickXorHash   *string `json:"QuickXorHash"`
	WebURL         string  `json:"webUrl"`
	CreatedAt      string  `json:"createdAt"`
	LastModifiedAt string  `json:"lastModifiedAt"`
}

type DriveItemsReferenceCache struct {
	APIRequestURL string          `json:"requestUrl"`
	Path          string          `json:"path"`
	ID            string          `json:"id"`
	DriveID       string          `json:"driveId"`
	DriveType     string          `json:"driveType"`
	LastUpdateAt  int64           `json:"lastUpdateAt"`
	Content       *DriveItemCache `json:"content"`
}

func (od *OneDrive) GraphAPIDriveItemsToDriveItemsCache(graphAPIDriveItems *GraphAPIDriveItems) (*DriveItemsCache, error) {
	if len(graphAPIDriveItems.DriveItems) == 0 {
		return nil, nil
	}

	driveItems := []DriveItemCache{}
	newDriveItem := DriveItemCache{}
	for _, driveItem := range graphAPIDriveItems.DriveItems {
		var childCount *int = nil
		if driveItem.Folder != nil {
			childCount = &driveItem.Folder.ChildCount
		}
		var mimeType *string = nil
		var quickXorHash *string = nil
		if driveItem.File != nil {
			mimeType = &driveItem.File.MimeType
			quickXorHash = &driveItem.File.Hashes.QuickXorHash
		}
		newDriveItem = DriveItemCache{
			ID:             driveItem.ID,
			DownloadURL:    driveItem.MicrosoftGraphDownloadURL,
			Name:           driveItem.Name,
			Size:           driveItem.Size,
			ChildCount:     childCount,
			MimeType:       mimeType,
			QuickXorHash:   quickXorHash,
			WebURL:         driveItem.WebURL,
			CreatedAt:      driveItem.CreatedDateTime,
			LastModifiedAt: driveItem.LastModifiedDateTime,
		}
		driveItems = append(driveItems, newDriveItem)
	}

	// "path": "/drive/root:/cdn/www.youtube.com/watch/2C4ry6O2enU/videoplayback/1080p60"
	item := graphAPIDriveItems.DriveItems[0].ParentReference
	path := item.Path
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += path
	reqURL += ":/children"

	driveItemsReference := DriveItemsReferenceCache{
		APIRequestURL: reqURL,
		Path:          path,
		ID:            item.ID,
		DriveID:       item.DriveID,
		DriveType:     item.DriveType,
		LastUpdateAt:  time.Now().Unix(),
	}

	driveItemsCache := &DriveItemsCache{
		DriveItems:          driveItems,
		DriveItemsReference: driveItemsReference,
	}

	log.Println("od.GraphAPIDriveItemsToDriveItemsCache " + path)
	return driveItemsCache, nil
}

func (od *OneDrive) GetGraphAPIDriveItemFromCache(path string) (*DriveItemsCache, error) {
	path, filename := RegularPathToPathFilename(path)
	path = od.PathToGraphAPIDriveItemsPath(path)
	for _, driveItemsCache := range od.DriveItemsCaches {
		driveItemsReference := driveItemsCache.DriveItemsReference
		if driveItemsReference.Path == path {
			if time.Now().Unix()-driveItemsReference.LastUpdateAt > 3600 {
				return nil, errors.New("CacheExpired")
			}
			driveItems := driveItemsCache.DriveItems
			for _, driveItem := range driveItems {
				if driveItem.Name == filename {
					driveItemsReference.Content = &driveItem
					newDriveItemsCache := DriveItemsCache{
						DriveItemsReference: driveItemsReference,
					}
					log.Println("HIT Cache "+path, driveItemsReference.LastUpdateAt)
					return &newDriveItemsCache, nil
				}
			}
		}
	}
	return nil, errors.New("NoCacheRecord")
}

func (od *OneDrive) GetGraphAPIDriveItemsFromCache(path string) (*DriveItemsCache, error) {
	path = od.PathToGraphAPIDriveItemsPath(path)
	for _, driveItemsCache := range od.DriveItemsCaches {
		driveItemsReference := driveItemsCache.DriveItemsReference
		if driveItemsReference.Path == path {
			if time.Now().Unix()-driveItemsReference.LastUpdateAt > 3600 {
				return nil, errors.New("CacheExpired")
			}
			log.Println("HIT Cache "+path, driveItemsReference.LastUpdateAt)
			return &driveItemsCache, nil
		}
	}
	return nil, errors.New("NoCacheRecord")
}

// DriveCache cache structure
type DriveCache struct {
	Path                 string           `json:"path"`
	RequestURL           string           `json:"requestUrl"`
	Name                 string           `json:"name"`
	Size                 int              `json:"size"`
	Items                []DriveCacheItem `json:"items"`
	CreatedDateTime      string           `json:"createdDateTime"`
	LastModifiedDateTime string           `json:"lastModifiedDateTime"`
	UpdateAt             int64            `json:"updateAt"`
}

type DriveCacheItem struct {
	Name                 string `json:"name"`
	Size                 int    `json:"size"`
	ChildCount           int    `json:"childCount"`
	DownloadURL          string `json:"downloadUrl"`
	MimeType             string `json:"mimeType"`
	QuickXorHash         string `json:"quickXorHash"`
	CreatedDateTime      string `json:"createdDateTime"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

type DriveCacheContentURL struct {
	RequestURL  string  `json:"requestUrl"`
	ResponseURL url.URL `json:"url"`
	UpdateAt    int64   `json:"updateAt"`
}

func (od *OneDrive) HitDriveItemsCaches(path string) (*DriveItemsCache, error) {
	driveItemsCache, err := od.GetGraphAPIDriveItemsFromCache(path)
	if err == nil {
		return driveItemsCache, nil
	}
	driveItemsCache, err = od.GetGraphAPIDriveItemFromCache(path)
	if err == nil {
		return driveItemsCache, nil
	}

	graphAPIDriveItems, err := od.GetGraphAPIDriveItems(path)
	if err != nil {
		return nil, err
	}
	driveItemsCache, err = od.GraphAPIDriveItemsToDriveItemsCache(graphAPIDriveItems)
	if err != nil {
		_ = driveItemsCache
		return nil, err
	}
	if driveItemsCache == nil {
		driveItemsCache, err = od.HitDriveItemCaches(path)
		if err == nil {
			return driveItemsCache, nil
		}
		driveItemsReference := DriveItemsReferenceCache{
			Path:         od.PathToGraphAPIDriveItemsPath(path),
			LastUpdateAt: time.Now().Unix(),
		}
		driveItemsCache = &DriveItemsCache{
			DriveItemsReference: driveItemsReference,
		}
	}
	od.DriveItemsCaches = append(od.DriveItemsCaches, *driveItemsCache)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	// go od.AutoCacheDrivePathContentURL(*driveCache)
	return driveItemsCache, nil
}

func (od *OneDrive) HitDriveItemCaches(path string) (*DriveItemsCache, error) {
	driveItemsCache, err := od.GetGraphAPIDriveItemFromCache(path)
	if err == nil {
		return driveItemsCache, nil
	}
	subPath, _ := RegularPathToPathFilename(path)

	graphAPIDriveItems, err := od.GetGraphAPIDriveItems(subPath)
	if err != nil {
		return nil, err
	}
	driveItemsCache, err = od.GraphAPIDriveItemsToDriveItemsCache(graphAPIDriveItems)
	if err != nil {
		_ = driveItemsCache
		return nil, err
	}
	if driveItemsCache == nil {
		return nil, errors.New("NoParentDriveItem")
	}
	od.DriveItemsCaches = append(od.DriveItemsCaches, *driveItemsCache)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	return od.GetGraphAPIDriveItemFromCache(path)
}

func (od *OneDrive) HitDriveItemContentURLCaches(path string) (*url.URL, error) {
	driveItemsCache, err := od.HitDriveItemCaches(path)
	if err != nil {
		return nil, err
	}
	return url.Parse(*driveItemsCache.DriveItemsReference.Content.DownloadURL)
}

func (od *OneDrive) CacheDrivePathContentURL(path string) (*url.URL, error) {
	reqURL := od.DrivePathContentToURL(path)

	for i, item := range od.DriveCacheContentURL {
		if item.RequestURL == reqURL {
			if time.Now().Unix()-item.UpdateAt > 1200 {
				log.Println("Updating " + item.RequestURL)
				driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
				if err != nil {
					return nil, err
				}
				od.DriveCacheContentURL[i] = *driveCacheContentURL
				return &driveCacheContentURL.ResponseURL, nil
			}
			log.Println("HIT Cache "+item.RequestURL, item.UpdateAt)
			return &item.ResponseURL, nil
		}
	}

	driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
	if err != nil {
		return nil, err
	}

	od.DriveCacheContentURL = append(od.DriveCacheContentURL, *driveCacheContentURL)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	return &driveCacheContentURL.ResponseURL, nil
}

func (od *OneDrive) AutoCacheDrivePathContentURL(driveCache DriveCache) error {
	for _, pItem := range driveCache.Items {
		if pItem.DownloadURL != "" {
			reqURL := od.DrivePathContentToURL(pItem.DownloadURL)
			func() {
				for i, item := range od.DriveCacheContentURL {
					if item.RequestURL == reqURL {
						if time.Now().Unix()-item.UpdateAt > 1200 {
							log.Println("Updating " + item.RequestURL)
							driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
							if err != nil {
								log.Println(err)
								return
							}
							od.DriveCacheContentURL[i] = *driveCacheContentURL
						}
						log.Println("HIT Cache "+item.RequestURL, item.UpdateAt)
					}
				}

				driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
				if err != nil {
					log.Println(err)
					return
				}

				od.DriveCacheContentURL = append(od.DriveCacheContentURL, *driveCacheContentURL)
			}()
		}
	}

	if err := od.SaveDriveCacheFile(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (od *OneDrive) getDriveCacheContentURL(reqURL string) (*DriveCacheContentURL, error) {
	log.Println("GET " + reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+od.MicrosoftGraphAPIToken.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	driveCacheContentURL := DriveCacheContentURL{
		RequestURL:  reqURL,
		ResponseURL: *resp.Request.URL,
		UpdateAt:    time.Now().Unix(),
	}

	return &driveCacheContentURL, nil
}
