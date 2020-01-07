package core

import (
	"encoding/json"
	"errors"
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
	DriveItemsCaches       []DriveItemsCache       `json:"driveItemsCaches"`
	DriveCache             []DriveCache            `json:"driveCache"`
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
			log.Println(od.AppRegistrationConfig.EndPointURI + "/authorize?client_id=" + clientID + "&scope=" + grantScope + "&response_type=code&redirect_uri=" + redirectURI)
		}
		// return nil, errors.New("Invalid Microsoft Graph API Token Grant Type")
		return nil, nil
	}
	data.Set("client_id", od.AppRegistrationConfig.ClientID)
	data.Set("client_secret", od.AppRegistrationConfig.ClientSecret)
	data.Set("redirect_uri", od.AppRegistrationConfig.RedirectURIs[0])
	return strings.NewReader(data.Encode()), nil
}

func (od *OneDrive) GetMicrosoftGraphAPIToken() error {
	endPointURI := od.AppRegistrationConfig.EndPointURI
	endPointURI += "/token"
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

	if err = json.Unmarshal(body, &od.MicrosoftGraphAPIToken); err != nil {
		return err
	}
	if od.MicrosoftGraphAPIToken == nil {
		log.Println(string(body))
		return errors.New("GetMicrosoftGraphAPITokenRequestError")
	}

	od.DriveDescriptionConfig.RefreshToken = od.MicrosoftGraphAPIToken.RefreshToken

	return nil
}

func RegularPathToPathFilename(rPath string) (path, filename string) {
	strS := strings.Split(rPath, "/")
	strR := ""
	for _, str := range strS {
		if str != "" {
			strR += "/" + str
		}
	}
	if strR == "" {
		strR += "/"
	}
	strRR := strings.Split(strR, "/")
	n := len(strRR)
	if n > 1 {
		for _, str := range strRR[0 : n-1] {
			if str != "" {
				path += "/" + str
			}
		}
		filename = strRR[n-1]
	}
	return path, filename
}

func RegularPath(path string) (rPath string) {
	// any path to "/" or "/path/to"
	path, _ = url.QueryUnescape(path)
	strD := strings.Split(path, "#")
	strQ := strings.Split(strD[0], "?")
	pathRaw := strQ[0]
	pathQuery := ""
	for i, str := range strQ {
		if i != 0 {
			pathQuery += str
		}
	}
	strS := strings.Split(pathRaw, "/")
	for _, str := range strS {
		if str != "" {
			rPath += "/" + str
		}
	}
	if rPath == "" {
		rPath += "/"
	}
	if pathQuery != "" {
		rPath += "?" + pathQuery
	}
	return rPath
}

func RegularRootPath(path string) (str string) {
	length := len(path)
	if length > 0 {
		path = RegularPath(path)
		if path == "/" || path == "/root" || path == "/root:" {
			str = "/drive/root:"
		} else {
			str = "/drive/root:" + path
		}
	} else {
		str = "/drive/root:"
	}
	return str
}

func (od *OneDrive) DrivePathToURL(path string) string {
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += "/me"
	reqURL += RegularRootPath(od.DriveDescriptionConfig.RootPath)
	reqURL += path
	reqURL += "?expand=children($select=name,size,file,folder,parentReference,createdDateTime,lastModifiedDateTime)"
	return reqURL
}

func (od *OneDrive) DrivePathContentToURL(path string) string {
	reqURL := od.DriveDescriptionConfig.EndPointURI
	reqURL += "/me"
	reqURL += RegularRootPath(od.DriveDescriptionConfig.RootPath)
	path = RegularPath(path)
	if path != "/" {
		reqURL += path
	}
	reqURL += ":/content"
	return reqURL
}

// DriveItemsPayload presents drive items cache structure
type DriveItemsPayload struct {
	DriveItems          []DriveItemPayload         `json:"driveItem"`
	DriveItemsReference DriveItemsReferencePayload `json:"driveItemReference"`
}

type DriveItemPayload struct {
	Path           string  `json:"path"`
	Name           string  `json:"name"`
	Size           int     `json:"size"`
	ChildCount     *int    `json:"childCount"`
	MimeType       *string `json:"mimeType"`
	QuickXorHash   *string `json:"QuickXorHash"`
	CreatedAt      string  `json:"createdAt"`
	LastModifiedAt string  `json:"lastModifiedAt"`
}

type DriveItemsReferencePayload struct {
	Path         string            `json:"path"`
	LastUpdateAt int64             `json:"lastUpdateAt"`
	Content      *DriveItemPayload `json:"content"`
}

func (od *OneDrive) GetDriveItemsFromRootPath() (*DriveItemsPayload, error) {
	return od.GetDriveItemsFromPath("")
}

func (od *OneDrive) GetDriveItemsFromPath(path string) (*DriveItemsPayload, error) {
	driveItemsCache, err := od.HitDriveItemsCaches(path)
	if err != nil {
		return nil, err
	}

	driveItemsReferenceCache := driveItemsCache.DriveItemsReference
	rPath := od.GraphAPIDriveItemsPathToPath(driveItemsReferenceCache.Path)
	driveItemsReferencePayload := DriveItemsReferencePayload{
		Path:         rPath,
		LastUpdateAt: driveItemsReferenceCache.LastUpdateAt,
	}
	var driveItemPayloads []DriveItemPayload = nil
	if rPath != "/" {
		rPath += "/"
	}
	if driveItemsReferenceCache.Content != nil {
		content := driveItemsReferenceCache.Content
		driveItemsReferencePayload.Content = &DriveItemPayload{
			Path:           rPath + content.Name,
			Name:           content.Name,
			Size:           content.Size,
			ChildCount:     content.ChildCount,
			MimeType:       content.MimeType,
			QuickXorHash:   content.QuickXorHash,
			CreatedAt:      content.CreatedAt,
			LastModifiedAt: content.LastModifiedAt,
		}
	} else {
		driveItems := driveItemsCache.DriveItems
		for _, driveItem := range driveItems {
			newDriveItemPayload := DriveItemPayload{
				Path:           rPath + driveItem.Name,
				Name:           driveItem.Name,
				Size:           driveItem.Size,
				ChildCount:     driveItem.ChildCount,
				MimeType:       driveItem.MimeType,
				QuickXorHash:   driveItem.QuickXorHash,
				CreatedAt:      driveItem.CreatedAt,
				LastModifiedAt: driveItem.LastModifiedAt,
			}
			driveItemPayloads = append(driveItemPayloads, newDriveItemPayload)
		}
	}
	driveItemsPayload := DriveItemsPayload{
		DriveItems:          driveItemPayloads,
		DriveItemsReference: driveItemsReferencePayload,
	}

	return &driveItemsPayload, nil
}

func (od *OneDrive) GetDriveItemContentURLFromPath(path string) (*url.URL, error) {
	return od.HitDriveItemContentURLCaches(path)
}
