package core

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// OneDrive presends a single OneDrive client end-point.
type OneDrive struct {
	AppRegistrationConfig  AppRegistrationConfig   `json:"appRegistrationConfig"`
	DriveDescriptionConfig DriveDescriptionConfig  `json:"driveDescriptionConfig"`
	MicrosoftGraphAPIToken *MicrosoftGraphAPIToken `json:"microsoftGraphApiToken"`
	DriveCache             []DriveCache           `json:"driveCache"`
	DriveCacheContentURL   []DriveCacheContentURL  `json:"driveCacheContentUrl"`
}

// AppRegistrationConfig configures the app in Azure Active Directory admin center.
type AppRegistrationConfig struct {
	EndPointURI  string   `json:"endPointUri"`
	DisplayName  *string  `json:"displayName"`
	ClientID     string   `json:"clientId"`
	TenantID     *string  `json:"tenantId"`
	ObjectID     *string  `json:"objectId"`
	RedirectURIs []string `json:"redirectUris"`
	LogoutURL    *string  `json:"logoutUrl"`
	ClientSecret string   `json:"clientSecret"`
}

// DriveDescriptionConfig configures the OneDrive local client.
type DriveDescriptionConfig struct {
	EndPointURI      string            `json:"endPointUri"`
	GrantScope       string            `json:"grantScope"`
	RootPath         string            `json:"rootPath"`
	Code             string            `json:"code"`
	RefreshToken     string            `json:"refreshToken"`
	VolumeMounts     *[]VolumeMount    `json:"volumeMounts"`
	DriveCacheConfig *DriveCacheConfig `json:"driveCacheConfig"`
}

// VolumeMount configures the volume mounts.
type VolumeMount struct {
	SourcePath string `json:"sourcePath"`
	MountPath  string `json:"mountPath"`
}

// DriveCacheConfig configures the drive files cache.
type DriveCacheConfig struct {
	CacheEabled           bool      `json:"cacheEabled"`
	CacheList             *[]string `json:"cacheList"`
	FileRefreshInterval   int       `json:"fileRefreshInterval"`
	FolderRefreshInterval int       `json:"folderRefreshInterval"`
}

// MicrosoftGraphAPIToken shows the Microsoft Graph API Token structure.
type MicrosoftGraphAPIToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	CreatedDateTime      string               `json:"createdDateTime"`
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

func (od *OneDrive) Run() error {
	if err := od.GetMicrosoftGraphAPIToken(); err != nil {
		return err
	}
	if err := od.SaveConfigFile(); err != nil {
		return err
	}
	if err := od.LoadDriveCacheFile(); err != nil {
		return err
	}
	if err := od.Cron(); err != nil {
		return err
	}
	return nil
}

func (od *OneDrive) getMicrosoftGraphAPITokenPostForm() (io.Reader, error) {
	data := url.Values{}
	if od.DriveDescriptionConfig.RefreshToken != "" {
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", od.DriveDescriptionConfig.RefreshToken)
	} else if od.DriveDescriptionConfig.Code != "" {
		data.Set("grant_type", "authorization_code")
		data.Set("code", od.DriveDescriptionConfig.Code)
	} else {
		log.Println("Invalid Microsoft Graph API Token Grant Type, use the following URLs to GET code")
		clientID := od.AppRegistrationConfig.ClientID
		grantScope := url.QueryEscape(od.DriveDescriptionConfig.GrantScope)
		for _, redirectURI := range od.AppRegistrationConfig.RedirectURIs {
			log.Println("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + clientID + "&scope=" + grantScope + "&response_type=code&redirect_uri=" + redirectURI)
		}
		// return nil, errors.New("Invalid Microsoft Graph API Token Grant Type")
		return nil, nil
	}
	data.Set("client_id", od.AppRegistrationConfig.ClientID)
	data.Set("client_secret", od.AppRegistrationConfig.ClientSecret)
	data.Set("redirect_uri", od.AppRegistrationConfig.RedirectURIs[3])
	return strings.NewReader(data.Encode()), nil
}

func (od *OneDrive) GetMicrosoftGraphAPIToken() error {
	endPointURI := od.AppRegistrationConfig.EndPointURI
	postForm, err := od.getMicrosoftGraphAPITokenPostForm()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", endPointURI, postForm)
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
	if err = json.Unmarshal([]byte(body), &od.MicrosoftGraphAPIToken); err != nil {
		return err
	}
	od.DriveDescriptionConfig.RefreshToken = od.MicrosoftGraphAPIToken.RefreshToken

	return nil
}

func RegularPath(path string) (str string) {
	// any path to "/" or "/path/to"
	path, _ = url.QueryUnescape(path)
	length := len(path)
	if length > 0 {
		if path == "/" {
			str = "/"
		} else if path[0] == '/' && path[length-1] == '/' {
			str = path[0:(length - 1)]
		} else if path[0] == '/' && path[length-1] != '/' {
			str = path
		} else if path[0] != '/' && path[length-1] == '/' {
			str = "/" + path[0:(length-1)]
		} else if path[0] != '/' && path[length-1] != '/' {
			str = "/" + path
		}
	} else {
		str = "/"
	}
	return str
}

func RegularRootPath(path string) (str string) {
	length := len(path)
	if length > 0 {
		path = RegularPath(path)
		if path == "/" || path == "/root" || path == "/root:" {
			str = "/root:"
		} else {
			str = "/root:" + path
		}
	} else {
		str = "/root:"
	}
	return str
}

func (od *OneDrive) DrivePathToURL(path string) string {
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += RegularRootPath(od.DriveDescriptionConfig.RootPath)
	reqURL += path
	reqURL += "?expand=children($select=name,size,file,folder,parentReference,createdDateTime,lastModifiedDateTime)"
	return reqURL
}

func (od *OneDrive) GetDriveRootPath() (*DriveCache, error) {
	return od.GetDrivePath("")
}

func (od *OneDrive) GetDrivePath(path string) (*DriveCache, error) {
	return od.CacheDrivePath(path)
}

func (od *OneDrive) DrivePathContentToURL(path string) string {
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += RegularRootPath(od.DriveDescriptionConfig.RootPath)
	reqURL += RegularPath(path)
	reqURL += ":/content"
	return reqURL
}

func (od *OneDrive) GetDrivePathContentURL(path string) (*url.URL, error) {
	return od.CacheDrivePathContentURL(path)
}
