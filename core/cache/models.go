package cache

import "github.com/AirWSW/onedrive/graphapi"

type DriveCacheCollection struct {
	
	MicrosoftGraphDriveItemCache []MicrosoftGraphDriveItemCache `json:"microsoftGraphDriveItemCache,omitempty"`
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
	Status       string `json:"status"` // Wait, Caching, Cached, Failed
}
