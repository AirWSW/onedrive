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
	Path     string    `json:"path"`
	URL      url.URL   `json:"url"`
	UpdateAt time.Time `json:"updateAt"`
}

func (od *OneDrive) CacheDrivePath(path string) (*DriveCache, error) {
	path = RegularPath(path)
	reqURL := od.DrivePathToURL(path)

	for _, driveCache := range od.DriveCache {
		if driveCache.Path == reqURL {
			log.Println("Read drive cache timestamp: ", driveCache.UpdateAt)
			return &driveCache, nil
		}
	}

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
		Path:                 reqURL,
		Name:                 driveFolder.Name,
		Size:                 driveFolder.Size,
		CreatedDateTime:      driveFolder.CreatedDateTime,
		LastModifiedDateTime: driveFolder.LastModifiedDateTime,
		Items:                driveCacheItem,
		UpdateAt:             time.Now(),
	}

	od.DriveCache = append(od.DriveCache, driveCache)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	go od.AutoCacheDrivePathContentURL(driveCache)

	return &driveCache, nil
}

func (od *OneDrive) CacheDrivePathContentURL(path string) (*url.URL, error) {
	reqURL := od.DrivePathContentToURL(path)

	for _, driveCacheContentURL := range od.DriveCacheContentURL {
		if driveCacheContentURL.Path == reqURL {
			log.Println("Read content cache timestamp: ", driveCacheContentURL.UpdateAt)
			return &driveCacheContentURL.URL, nil
		}
	}

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
		Path:     reqURL,
		URL:      *resp.Request.URL,
		UpdateAt: time.Now(),
	}
	od.DriveCacheContentURL = append(od.DriveCacheContentURL, driveCacheContentURL)

	if err := od.SaveDriveCacheFile(); err != nil {
		return nil, err
	}

	return &driveCacheContentURL.URL, nil
}

func (od *OneDrive) AutoCacheDrivePathContentURL(driveCache DriveCache) error {
	for _, item := range driveCache.Items {
		if item.DownloadURL != "" {
			reqURL := od.DrivePathContentToURL(item.DownloadURL)
			log.Println(reqURL)

			func() {
				for _, driveCacheContentURL := range od.DriveCacheContentURL {
					if driveCacheContentURL.Path == reqURL {
						log.Println("Read content cache timestamp: ", driveCacheContentURL.UpdateAt)
						return
					}
				}

				req, err := http.NewRequest("GET", reqURL, nil)
				if err != nil {
					log.Println(err)
					return
				}
				req.Header.Add("Authorization", "Bearer "+od.MicrosoftGraphAPIToken.AccessToken)
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					return
				}

				driveCacheContentURL := DriveCacheContentURL{
					Path:     reqURL,
					URL:      *resp.Request.URL,
					UpdateAt: time.Now(),
				}
				od.DriveCacheContentURL = append(od.DriveCacheContentURL, driveCacheContentURL)
			}()
		}
	}

	if err := od.SaveDriveCacheFile(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
