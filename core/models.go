package core

import "github.com/AirWSW/onedrive/graphapi"

// OneDriveCollection collects all OneDrives
type OneDriveCollection struct {
	OneDrives []*OneDrive `json:"oneDrives"`
}

// OneDrive describes a OneDrive
type OneDrive struct {
	MicrosoftEndPoints     graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
	AzureADAppRegistration graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
	AzureADAuthFlowContext graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
	OneDriveDescription    OneDriveDescription             `json:"oneDriveDescription"`
	MicrosoftGraphAPI      *graphapi.MicrosoftGraphAPI     `json:"microsoftGraphApi"`
}

// OneDriveDescription describes the OneDrive local client
type OneDriveDescription struct {
	RootPath         string                       `json:"rootPath"`
	VolumeMounts     []*VolumeMount               `json:"volumeMounts"`
	CacheConfig      *DriveCacheConfig            `json:"driveCacheConfig"`
	DriveDescription *graphapi.MicrosoftGraphDrive `json:"driveDescription"`
}
