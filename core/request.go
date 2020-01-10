package core

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/AirWSW/onedrive/graphapi"
)

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveRaw(str string) ([]byte, error) {
	return od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(str)
}

func (od *OneDrive) PostMicrosoftGraphAPIMeDriveRaw(str string, postBody io.Reader) ([]byte, error) {
	return od.MicrosoftGraphAPI.UseMicrosoftGraphAPIPost(str, postBody)
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDrive() error {
	odd := od.OneDriveDescription
	bytes, err := od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDrivePath(""))
	if err != nil {
		return err
	}
	microsoftGraphDrive := graphapi.MicrosoftGraphDrive{}
	if err := json.Unmarshal(bytes, &microsoftGraphDrive); err != nil {
		return err
	}
	return nil
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveItem(str string) (*graphapi.MicrosoftGraphDriveItem, error) {
	odd := od.OneDriveDescription
	reqURL := odd.UseMicrosoftGraphAPIMeDriveItem(str)
	strURL, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if strURL.Scheme == "https" {
		reqURL = str
	}
	bytes, err := od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(reqURL)
	if err != nil {
		return nil, err
	}
	microsoftGraphDriveItem := graphapi.MicrosoftGraphDriveItem{}
	if err := json.Unmarshal(bytes, &microsoftGraphDriveItem); err != nil {
		return nil, err
	}
	return &microsoftGraphDriveItem, nil
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveChildren(str string) (*graphapi.MicrosoftGraphDriveItemCollection, error) {
	odd := od.OneDriveDescription
	bytes := []byte{}

	url, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "https" {
		bytes, err = od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveChildren(str))
		if err != nil {
			return nil, err
		}
	} else {
		bytes, err = od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(str)
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

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveExpandChildren(str string) error {
	odd := od.OneDriveDescription
	bytes, err := od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveExpandChildrenPath(str))
	if err != nil {
		return err
	}
	microsoftGraphDriveItem := graphapi.MicrosoftGraphDriveItem{}
	if err := json.Unmarshal(bytes, &microsoftGraphDriveItem); err != nil {
		return err
	}
	return nil
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveContent(str string) ([]byte, error) {
	odd := od.OneDriveDescription
	return od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet(odd.UseMicrosoftGraphAPIMeDriveContentPath(str))
}

func (od *OneDrive) GetMicrosoftGraphAPIMeDriveChildrenRequest(str string) (*MicrosoftGraphDriveItemCache, error) {
	microsoftGraphDriveItemCache := &MicrosoftGraphDriveItemCache{
		Children: []MicrosoftGraphDriveItemCache{},
	}

	url, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "https" {
		microsoftGraphDriveItem, err := od.GetMicrosoftGraphAPIMeDriveItem(str)
		if err != nil {
			return nil, err
		}
		microsoftGraphDriveItemCache, err = od.DriveItemToCache(microsoftGraphDriveItem)
		if err != nil {
			return nil, err
		}
	}

	microsoftGraphDriveItemCollection, err := od.GetMicrosoftGraphAPIMeDriveChildren(str)
	if err != nil {
		return nil, err
	}
	for _, value := range microsoftGraphDriveItemCollection.Value {
		newMicrosoftGraphDriveItemCache := &MicrosoftGraphDriveItemCache{
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
		newMicrosoftGraphDriveItemCache, err := od.GetMicrosoftGraphAPIMeDriveChildrenRequest(*microsoftGraphDriveItemCollection.AtODataNextLink)
		if err != nil {
			return newMicrosoftGraphDriveItemCache, err
		}
		for _, children := range newMicrosoftGraphDriveItemCache.Children {
			microsoftGraphDriveItemCache.Children = append(microsoftGraphDriveItemCache.Children, children)
		}
	}

	return microsoftGraphDriveItemCache, nil
}
