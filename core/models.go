package core

import (
	"time"

	"github.com/AirWSW/onedrive/core/api"
	"github.com/AirWSW/onedrive/core/cache"
	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/core/upload"
	"github.com/AirWSW/onedrive/graphapi"
)

// OneDrive describes a OneDrive
type OneDrive struct {
	MicrosoftEndPoints     graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
	AzureADAppRegistration graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
	OneDriveDescription    description.OneDriveDescription `json:"oneDriveDescription"`
	MicrosoftGraphAPI      api.MicrosoftGraphAPI           `json:"microsoftGraphApi,omitempty"`
	DriveCacheCollection   cache.DriveCacheCollection      `json:"driveCacheCollection,omitempty"`
	UploaderCollection     upload.UploaderCollection       `json:"uploaderCollection,omitempty"`
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
