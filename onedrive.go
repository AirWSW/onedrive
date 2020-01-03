package onedrive

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OneDrive struct {
	DriveConfig     DriveConfig     `json:"drive_config"`
	RegisterConfig  RegisterConfig  `json:"register_config"`
	MSGraphAPIToken MSGraphAPIToken `json:"ms_graph_api_token"`
	DriveCache      []DriveCache    `json:"drive_cache"`
}

type RegisterConfig struct {
	EndPointURI  string `json:"end_point_uri"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type DriveConfig struct {
	EndPointURI           string `json:"end_point_uri"`
	RootPath              string `json:"root_path"`
	FileRefreshInterval   int    `json:"file_refresh_interval"`
	FolderRefreshInterval int    `json:"folder_refresh_interval"`
	Code                  string `json:"code"`
	RefreshToken          string `json:"refresh_token"`
}

type MSGraphAPIToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Drive cache structure

type DriveCache struct {
	Path                 string           `json:"path"`
	Name                 string           `json:"name"`
	Size                 int              `json:"size"`
	Items                []DriveCacheItem `json:"items"`
	CreatedDateTime      string           `json:"createdDateTime"`
	LastModifiedDateTime string           `json:"lastModifiedDateTime"`
	CacheLastUpdateAt    string           `json:"lastUpdateAt"`
}

type DriveCacheItem struct {
	Name                 string `json:"name"`
	Size                 int    `json:"size"`
	ChildCount           int    `json:"childCount"`
	DownloadURL          string `json:"downloadUrl"`
	MimeType             string `json:"mimeType"`
	QuickXorHash         string `json:"quickXorHash"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

// Drive structure

type DriveFolder struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	Size                 int                  `json:"size"`
	WebURL               string               `json:"webUrl"`
	CreatedDateTime      string               `json:"createdDateTime"`
	LastModifiedDateTime string               `json:"lastModifiedDateTime"`
	Folder               DriveChildrenFolder  `json:"folder"`
	ParentReference      DriveParentReference `json:"parentReference"`
	Children             []DriveChildrenItem  `json:"children"`
}

type DriveChildrenItem struct {
	Name                 string               `json:"name"`
	Size                 int                  `json:"size"`
	LastModifiedDateTime string               `json:"lastModifiedDateTime"`
	DownloadURL          string               `json:"@microsoft.graph.downloadUrl"`
	File                 DriveChildrenFile    `json:"file"`
	Folder               DriveChildrenFolder  `json:"folder"`
	ParentReference      DriveParentReference `json:"parentReference"`
}

type DriveChildrenFolder struct {
	ChildCount int                     `json:"childCount"`
	View       DriveChildrenFolderView `json:"view"`
}

type DriveChildrenFolderView struct {
	ViewType  string `json:"viewType"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type DriveChildrenFile struct {
	MimeType string                  `json:"mimeType"`
	Hashes   DriveChildrenFileHashes `json:"hashes"`
}

type DriveChildrenFileHashes struct {
	QuickXorHash string `json:"quickXorHash"`
	Sha1Hash     string `json:"sha1Hash"`
}

type DriveParentReference struct {
	ID        string `json:"id"`
	DriveID   string `json:"driveId"`
	DriveType string `json:"driveType"`
	Name      string `json:"name"`
	Path      string `json:"path"`
}

func CreateDefaultOneDrive() (od OneDrive, err error) {
	od.RegisterConfig = RegisterConfig{
		EndPointURI:  "https://login.chinacloudapi.cn/common/oauth2/v2.0/token", // "https://login.microsoftonline.com/common/oauth2/v2.0/token"
		Scope:        "user.read files.readwrite.all offline_access",
		ClientID:     "",
		ClientSecret: "",
		RedirectURI:  "http://localhost:8081/auth-redirect",
	}
	od.DriveConfig = DriveConfig{
		EndPointURI:           "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive", //"https://graph.microsoft.com/v1.0/me/drive"
		RootPath:              "root",
		FileRefreshInterval:   1200,
		FolderRefreshInterval: 600,
		Code:                  "",
		RefreshToken:          "",
	}
	return od, nil
}

func (od *OneDrive) getMSGraphAPITokenInput() io.Reader {
	data := url.Values{}
	if od.DriveConfig.RefreshToken != "" {
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", od.DriveConfig.RefreshToken)
	} else {
		data.Set("grant_type", "authorization_code")
		data.Set("code", od.DriveConfig.Code)
	}
	data.Set("client_id", od.RegisterConfig.ClientID)
	data.Set("client_secret", od.RegisterConfig.ClientSecret)
	data.Set("redirect_uri", od.RegisterConfig.RedirectURI)
	return strings.NewReader(data.Encode())
}

func (od *OneDrive) GetMSGraphAPIToken() error {
	url := od.RegisterConfig.EndPointURI
	input := od.getMSGraphAPITokenInput()
	req, err := http.NewRequest("POST", url, input)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(body), &od.MSGraphAPIToken); err != nil {
		return err
	}
	log.Println(string(body))
	return nil
}

func RegularPath(path string) (str string) {
	path, _ = url.QueryUnescape(path)
	if path != "" && path[0] == '/' {
		str = path[1:]
	} else if path != "" && path[0] != '/' {
		str = path
	} else {
		str = ""
	}
	n := len(str)
	if str[n-1] == '/' {
		str = str[:n-1]
	}
	return str
}

func (od *OneDrive) DrivePathToURL(path string) string {
	url := od.DriveConfig.EndPointURI
	rootPath := od.DriveConfig.RootPath
	if rootPath == "root" {
		url += "/root:"
	} else if rootPath[0] == '/' {
		url += "/root:" + rootPath
	}
	if path != "" && path[0] == '/' {
		url += path
	} else if path != "" && path[0] != '/' {
		url += "/" + path
	}
	url += "?expand=children($select=name,size,file,folder,parentReference,lastModifiedDateTime)"
	return url
}

func (od *OneDrive) CacheDrivePath(path string) (*DriveCache, error) {
	path = RegularPath(path)

	reqUrl := od.DrivePathToURL(path)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+od.MSGraphAPIToken.AccessToken)
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

	log.Println(string(body))
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
			LastModifiedDateTime: item.LastModifiedDateTime,
		})
	}
	driveCache := &DriveCache{
		Path:                 path,
		Name:                 driveFolder.Name,
		Size:                 driveFolder.Size,
		CreatedDateTime:      driveFolder.CreatedDateTime,
		LastModifiedDateTime: driveFolder.LastModifiedDateTime,
		Items:                driveCacheItem,
		CacheLastUpdateAt:    time.Now().String(),
	}
	return driveCache, nil
}

func (od *OneDrive) GetDriveRootPath() (*DriveCache, error) {
	return od.GetDrivePath("")
}

func (od *OneDrive) GetDrivePath(path string) (*DriveCache, error) {
	return od.CacheDrivePath(path)
}

func (od *OneDrive) DrivePathContentToURL(path string) string {
	url := od.DriveConfig.EndPointURI
	rootPath := od.DriveConfig.RootPath
	if rootPath == "root" {
		url += "/root:"
	} else if rootPath[0] == '/' {
		url += "/root:" + rootPath
	}
	if path != "" && path[0] == '/' {
		url += path
	} else if path != "" && path[0] != '/' {
		url += "/" + path
	}
	url += ":/content"
	return url
}

func (od *OneDrive) GetDrivePathContentURL(path string) (*url.URL, error) {
	url := od.DrivePathContentToURL(path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+od.MSGraphAPIToken.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Println(resp.Request.URL)
	return resp.Request.URL, nil
}
