package core

import (
	"time"

	"github.com/AirWSW/onedrive/graphapi"
)

// OneDriveCollection collects all OneDrives
type OneDriveCollection struct {
	OneDrives []*OneDrive `json:"oneDrives"`
}

// OneDrive describes a OneDrive
type OneDrive struct {
	MicrosoftEndPoints           graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
	AzureADAppRegistration       graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext       graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
	OneDriveDescription          OneDriveDescription             `json:"oneDriveDescription"`
	MicrosoftGraphAPI            *graphapi.MicrosoftGraphAPI     `json:"microsoftGraphApi,omitempty"`
	MicrosoftGraphDriveItemCache []MicrosoftGraphDriveItemCache  `json:"microsoftGraphDriveItemCache,omitempty"`
}

// OneDriveDescription describes the OneDrive local client
type OneDriveDescription struct {
	RootPath          string                        `json:"rootPath"`
	RefreshInterval   int64                         `json:"refreshInterval"`
	DriveVolumeMounts []DriveVolumeMount            `json:"driveVolumeMounts,omitempty"`
	CacheConfig       *DriveCacheConfig             `json:"driveCacheConfig,omitempty"`
	DriveDescription  *graphapi.MicrosoftGraphDrive `json:"driveDescription,omitempty"`
}

// DriveVolumeMount configures the volume mounts.
type DriveVolumeMount struct {
	Type     *string `json:"type"`
	Source   *string `json:"source"`
	Target   *string `json:"target"`
	Password *string `json:"password"`
}

// DriveCacheConfig configures the drive files cache.
type DriveCacheConfig struct {
	CacheEabled           bool      `json:"cacheEabled"`
	CacheList             *[]string `json:"cacheList"`
	FileRefreshInterval   int       `json:"fileRefreshInterval"`
	FolderRefreshInterval int       `json:"folderRefreshInterval"`
}

// MicrosoftGraphDriveItemCache describes the MicrosoftGraphDriveItem cache structure
type MicrosoftGraphDriveItemCache struct {
	CacheDescription *CacheDescription `json:"cacheDescription,omitempty"`

	CTag        string                         `json:"cTag"` // etag
	Description *string                        `json:"description,omitempty"`
	File        *graphapi.MicrosoftGraphFile   `json:"file,omitempty"`
	Folder      *graphapi.MicrosoftGraphFolder `json:"folder,omitempty"`
	Size        int64                          `json:"size"`

	/* relationships */
	Children []MicrosoftGraphDriveItemCache `json:"children,omitempty"`

	/* inherited from baseItem */
	ID              string                                `json:"id"` // identifier
	CreatedAt       int64                                 `json:"createdAt"`
	ETag            string                                `json:"eTag"`
	LastModifiedAt  int64                                 `json:"lastModifiedAt"`
	Name            string                                `json:"name"`
	ParentReference *graphapi.MicrosoftGraphItemReference `json:"parentReference,omitempty"`
	WebURL          string                                `json:"webUrl"`

	/* instance annotations */
	AtMicrosoftGraphDownloadURL string `json:"@microsoft.graph.downloadUrl"`
}

type CacheDescription struct {
	RequestURL   string `json:"requestUrl"`
	Path         string `json:"path"`
	LastUpdateAt int64  `json:"createdAt"`
	Status       string `json:"status"`
}

type DriveItemCachePayload struct {
	LastUpdateAt    time.Time                       `json:"lastUpdateAt"`
	Description     *string                         `json:"description,omitempty"`
	File            *graphapi.MicrosoftGraphFile    `json:"file,omitempty"`
	Folder          *graphapi.MicrosoftGraphFolder  `json:"folder,omitempty"`
	Size            int64                           `json:"size"`
	Children        []DriveItemCachePayload         `json:"children,omitempty"`
	CreatedAt       time.Time                       `json:"createdAt"`
	LastModifiedAt  time.Time                       `json:"lastModifiedAt"`
	Name            string                          `json:"name"`
	ParentReference *DriveItemCachePayloadReference `json:"parentReference,omitempty"`
	DownloadURL     *string                         `json:"downloadUrl,omitempty"`
}

type DriveItemCachePayloadReference struct {
	DriveType string `json:"driveType"` // personal, business, documentLibrary
	Path      string `json:"path"`
}
