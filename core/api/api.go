package api

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/AirWSW/onedrive/core/cache"
	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/graphapi"
)

type MicrosoftGraphAPI struct {
	graphapi.MicrosoftGraphAPI
}

// NewMicrosoftGraphAPI validates NewMicrosoftGraphAPIInput and assigns to api
func NewMicrosoftGraphAPI(input *graphapi.NewMicrosoftGraphAPIInput) (*MicrosoftGraphAPI, error) {
	graphapiMicrosoftGraphAPI, err := graphapi.NewMicrosoftGraphAPI(input)
	if err != nil {
		return nil, err
	}
	api := &MicrosoftGraphAPI{
		MicrosoftGraphAPI: *graphapiMicrosoftGraphAPI,
	}
	// return *MicrosoftGraphAPI as api
	return api, nil
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveRaw(str string) ([]byte, error) {
	return api.UseMicrosoftGraphAPIGet(str)
}

func (api *MicrosoftGraphAPI) PostMicrosoftGraphAPIMeDriveRaw(str string, postBody io.Reader) ([]byte, error) {
	return api.UseMicrosoftGraphAPIPost(str, postBody)
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDrive(odd *description.OneDriveDescription) error {
	bytes, err := api.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDrivePath(""))
	if err != nil {
		return err
	}
	microsoftGraphDrive := graphapi.MicrosoftGraphDrive{}
	if err := json.Unmarshal(bytes, &microsoftGraphDrive); err != nil {
		return err
	}
	return nil
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveItem(odd *description.OneDriveDescription, str string) (*graphapi.MicrosoftGraphDriveItem, error) {
	reqURL := odd.UseMicrosoftGraphAPIMeDriveItem(str)
	strURL, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if strURL.Scheme == "https" {
		reqURL = str
	}
	bytes, err := api.UseMicrosoftGraphAPIGet(reqURL)
	if err != nil {
		return nil, err
	}
	microsoftGraphDriveItem := graphapi.MicrosoftGraphDriveItem{}
	if err := json.Unmarshal(bytes, &microsoftGraphDriveItem); err != nil {
		return nil, err
	}
	return &microsoftGraphDriveItem, nil
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveChildren(odd *description.OneDriveDescription, str string) (*graphapi.MicrosoftGraphDriveItemCollection, error) {
	url, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	bytes := []byte{}
	if url.Scheme != "https" {
		bytes, err = api.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveChildren(str))
		if err != nil {
			return nil, err
		}
	} else {
		bytes, err = api.UseMicrosoftGraphAPIGet(str)
	}
	if err != nil {
		return nil, err
	}
	microsoftGraphDriveItemCollection := graphapi.MicrosoftGraphDriveItemCollection{}
	if err := json.Unmarshal(bytes, &microsoftGraphDriveItemCollection); err != nil {
		return nil, err
	}
	return &microsoftGraphDriveItemCollection, nil
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveExpandChildren(odd *description.OneDriveDescription, str string) error {
	bytes, err := api.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveExpandChildrenPath(str))
	if err != nil {
		return err
	}
	microsoftGraphDriveItem := graphapi.MicrosoftGraphDriveItem{}
	if err := json.Unmarshal(bytes, &microsoftGraphDriveItem); err != nil {
		return err
	}
	return nil
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveContent(odd *description.OneDriveDescription, str string) ([]byte, error) {
	return api.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveContentPath(str))
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIMeDriveChildrenRequest(odd *description.OneDriveDescription, str string) (*cache.MicrosoftGraphDriveItemCache, error) {
	microsoftGraphDriveItemCache := &cache.MicrosoftGraphDriveItemCache{
		Children: []cache.MicrosoftGraphDriveItemCache{},
	}

	url, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "https" {
		microsoftGraphDriveItem, err := api.GetMicrosoftGraphAPIMeDriveItem(odd, str)
		if err != nil {
			return nil, err
		}
		microsoftGraphDriveItemCache, err = cache.DriveItemToCache(microsoftGraphDriveItem)
		if err != nil {
			return nil, err
		}
	}

	microsoftGraphDriveItemCollection, err := api.GetMicrosoftGraphAPIMeDriveChildren(odd, str)
	if err != nil {
		return nil, err
	}
	for _, value := range microsoftGraphDriveItemCollection.Value {
		newMicrosoftGraphDriveItemCache := &cache.MicrosoftGraphDriveItemCache{
			CTag:                        value.CTag,
			Description:                 value.Description,
			File:                        value.File,
			Folder:                      value.Folder,
			Size:                        value.Size,
			ID:                          value.ID,
			CreatedAt:                   value.CreatedDateTime.Unix(),
			ETag:                        value.ETag,
			LastModifiedAt:              value.LastModifiedDateTime.Unix(),
			Name:                        value.Name,
			ParentReference:             value.ParentReference,
			WebURL:                      value.WebURL,
			AtMicrosoftGraphDownloadURL: value.AtMicrosoftGraphDownloadURL,
		}
		microsoftGraphDriveItemCache.Children = append(microsoftGraphDriveItemCache.Children, *newMicrosoftGraphDriveItemCache)
	}

	if microsoftGraphDriveItemCollection.AtODataNextLink != nil {
		newMicrosoftGraphDriveItemCache, err := api.GetMicrosoftGraphAPIMeDriveChildrenRequest(odd, *microsoftGraphDriveItemCollection.AtODataNextLink)
		if err != nil {
			return newMicrosoftGraphDriveItemCache, err
		}
		for _, children := range newMicrosoftGraphDriveItemCache.Children {
			microsoftGraphDriveItemCache.Children = append(microsoftGraphDriveItemCache.Children, children)
		}
	}

	return microsoftGraphDriveItemCache, nil
}

func (api *MicrosoftGraphAPI) UpdateMicrosoftGraphDriveItemCache(odd *description.OneDriveDescription, cacheDescription *cache.CacheDescription) (*cache.MicrosoftGraphDriveItemCache, error) {
	return api.GetMicrosoftGraphAPIMeDriveChildrenRequest(odd, cacheDescription.Path)
}
