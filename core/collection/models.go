package collection

import "github.com/AirWSW/onedrive/core"

// OneDriveCollection collects all OneDrives
type OneDriveCollection struct {
	IsDebugMode  *bool            `json:"isDebugMode"`
	PageTemplate *string          `json:"pageTemplate"`
	OneDrives    []*core.OneDrive `json:"oneDrives"`
}
