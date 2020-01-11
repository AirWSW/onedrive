package upload

import (
	"github.com/AirWSW/onedrive/graphapi"
)

type UploaderCollection struct {
	Uploaders []*Uploader `json:"uploaders"`
}

type Uploader struct {
	UploaderDescription *UploaderDescription `json:"uploaderDescription"`
	UploadSessions      []UploadSession      `json:"uploadSessions,omitempty"`
}

type UploaderDescription struct {
	UploaderReference    *UploaderReference                                    `json:"uploaderReference"`
	UploadableProperties *graphapi.MicrosoftGraphDriveItemUploadableProperties `json:"uploadableProperties"`
}

type UploaderReference struct {
	DriveType string  `json:"driveType"` // personal, business, documentLibrary
	Name      string  `json:"name,omitempty"`
	Size      int64   `json:"size"`
	Path      string  `json:"path"`
	UploadURL *string `json:"uploadUrl"`
}

type UploadSession struct {
	UploadSessionDescription *UploadSessionDescription             `json:"uploadSessionDescription"`
	UploadSessionReference   *graphapi.MicrosoftGraphUploadSession `json:"uploadSessionReference"`
}

// ContentLength: 26
type UploadSessionDescription struct {
	Status        string                               `json:"status"`
	ContentLength int64                                `json:"contentLength"`
	ContentRange  UploadSessionDescriptionContentRange `json:"contentRange"`
}

// ContentRange: bytes 0-25/128
type UploadSessionDescriptionContentRange struct {
	Type string `json:"type"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}
