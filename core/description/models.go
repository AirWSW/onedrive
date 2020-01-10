package description

import "github.com/AirWSW/onedrive/graphapi"

// OneDriveDescription describes the OneDrive local client
type OneDriveDescription struct {
	OneDriveName      string                        `json:"oneDriveName"`
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
