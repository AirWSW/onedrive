package upload

import (
	"github.com/AirWSW/onedrive/graphapi"
)

type UploadCollection struct {
	Uploads []*Upload `json:"uploads"`
}

type Upload struct {
	UploadDescription UploadDescription `json:"uploadDescription"`
	UploadSessions    []UploadSession   `json:"uploadSessions,omitempty"`
}

type UploadDescription struct {
	UploadableProperties graphapi.MicrosoftGraphDriveItemUploadableProperties `json:"uploadableProperties"`
}

type UploadSession struct {
	UploadSessionDescription UploadSessionDescription             `json:"uploadSessionDescription"`
	UploadSessionReference   graphapi.MicrosoftGraphUploadSession `json:"uploadSessionReference"`
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
