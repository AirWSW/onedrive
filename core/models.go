package core

import (
	"time"

	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/core/upload"
	"github.com/AirWSW/onedrive/graphapi"
)

// OneDrive describes a OneDrive
type OneDrive struct {
	MicrosoftEndPoints           graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
	AzureADAppRegistration       graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext       graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
	OneDriveDescription          description.OneDriveDescription `json:"oneDriveDescription"`
	MicrosoftGraphAPI            *graphapi.MicrosoftGraphAPI     `json:"microsoftGraphApi,omitempty"`
	MicrosoftGraphDriveItemCache []MicrosoftGraphDriveItemCache  `json:"microsoftGraphDriveItemCache,omitempty"`
	UploadCollection             *upload.UploadCollection        `json:"uploadCollection,omitempty"`
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
	Description    *string                         `json:"description,omitempty"`
	File           *graphapi.MicrosoftGraphFile    `json:"file,omitempty"`
	Folder         *graphapi.MicrosoftGraphFolder  `json:"folder,omitempty"`
	Size           int64                           `json:"size"`
	Children       []DriveItemCachePayload         `json:"children,omitempty"`
	CreatedAt      time.Time                       `json:"createdAt"`
	LastModifiedAt time.Time                       `json:"lastModifiedAt"`
	Name           string                          `json:"name"`
	Reference      *DriveItemCachePayloadReference `json:"reference,omitempty"`
	DownloadURL    *string                         `json:"downloadUrl,omitempty"`
}

type DriveItemCachePayloadReference struct {
	LastUpdateAt time.Time `json:"lastUpdateAt"`
	DriveType    string    `json:"driveType"` // personal, business, documentLibrary
	Path         string    `json:"path"`
}
