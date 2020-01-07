package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// DriveCache cache structure
type DriveCache struct {
	Path                 string           `json:"path"`
	RequestURL           string           `json:"requestUrl"`
	Name                 string           `json:"name"`
	Size                 int              `json:"size"`
	Items                []DriveCacheItem `json:"items"`
	CreatedDateTime      string           `json:"createdDateTime"`
	LastModifiedDateTime string           `json:"lastModifiedDateTime"`
	UpdateAt             time.Time        `json:"updateAt"`
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
	RequestURL  string    `json:"requestUrl"`
	ResponseURL url.URL   `json:"url"`
	UpdateAt    time.Time `json:"updateAt"`
}

func (od *OneDrive) CacheDrivePath(path string) (*DriveCache, error) {
	path = RegularPath(path)
	reqURL := od.DrivePathToURL(path)

	for i, driveCache := range od.DriveCache {
		if driveCache.Path == reqURL {
			if time.Now().Sub(driveCache.UpdateAt) < time.Second*1800 {
				driveCache, err := od.getDriveCache(path, reqURL)
				if err != nil {
					return nil, err
				}
				od.DriveCache[i] = *driveCache
				return driveCache, nil
			}
			log.Println("Read drive cache timestamp: ", driveCache.UpdateAt)
			return &driveCache, nil
		}
	}

	driveCache, err := od.getDriveCache(path, reqURL)
	if err != nil {
		return nil, err
	}

	od.DriveCache = append(od.DriveCache, *driveCache)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	go od.AutoCacheDrivePathContentURL(*driveCache)

	return driveCache, nil
}

func (od *OneDrive) getDriveCache(path, reqURL string) (*DriveCache, error) {
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
	driveFolder := DriveFolder{}
	if err = json.Unmarshal([]byte(body), &driveFolder); err != nil {
		return nil, err
	}

	driveCacheItem := []DriveCacheItem{}
	for _, item := range driveFolder.Children {
		downloadURL := ""
		if item.File.MimeType != "" {
			name, _ := url.QueryUnescape(item.Name)
			downloadURL += path + "/" + name
		}
		driveCacheItem = append(driveCacheItem, DriveCacheItem{
			Name:                 item.Name,
			Size:                 item.Size,
			ChildCount:           item.Folder.ChildCount,
			DownloadURL:          downloadURL,
			MimeType:             item.File.MimeType,
			QuickXorHash:         item.File.Hashes.QuickXorHash,
			CreatedDateTime:      item.CreatedDateTime,
			LastModifiedDateTime: item.LastModifiedDateTime,
		})
	}
	driveCache := DriveCache{
		RequestURL:           reqURL,
		Path:                 path,
		Name:                 driveFolder.Name,
		Size:                 driveFolder.Size,
		CreatedDateTime:      driveFolder.CreatedDateTime,
		LastModifiedDateTime: driveFolder.LastModifiedDateTime,
		Items:                driveCacheItem,
		UpdateAt:             time.Now(),
	}

	return &driveCache, nil
}

func (od *OneDrive) CacheDrivePathContentURL(path string) (*url.URL, error) {
	reqURL := od.DrivePathContentToURL(path)

	for i, driveCacheContentURL := range od.DriveCacheContentURL {
		if driveCacheContentURL.RequestURL == reqURL {
			if time.Now().Sub(driveCacheContentURL.UpdateAt) < time.Second*1800 {
				driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
				if err != nil {
					return nil, err
				}
				od.DriveCacheContentURL[i] = *driveCacheContentURL
				return &driveCacheContentURL.ResponseURL, nil
			}
			log.Println("Read content cache timestamp: ", driveCacheContentURL.UpdateAt)
			return &driveCacheContentURL.ResponseURL, nil
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
	for _, item := range driveCache.Items {
		if item.DownloadURL != "" {
			reqURL := od.DrivePathContentToURL(item.DownloadURL)
			log.Println(reqURL)

			func() {
				for i, driveCacheContentURL := range od.DriveCacheContentURL {
					if driveCacheContentURL.RequestURL == reqURL {
						if time.Now().Sub(driveCacheContentURL.UpdateAt) < time.Second*1800 {
							driveCacheContentURL, err := od.getDriveCacheContentURL(reqURL)
							if err != nil {
								log.Println(err)
								return
							}
							od.DriveCacheContentURL[i] = *driveCacheContentURL
						}
						log.Println("Read content cache timestamp: ", driveCacheContentURL.UpdateAt)
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

	driveCacheContentURL := DriveCacheContentURL{
		RequestURL:  reqURL,
		ResponseURL: *resp.Request.URL,
		UpdateAt:    time.Now(),
	}

	return &driveCacheContentURL, nil
}
