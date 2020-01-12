package description

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/AirWSW/onedrive/core/utils"
	"github.com/AirWSW/onedrive/graphapi"
)

type MicrosoftGraphAPI interface {
	UseMicrosoftGraphAPIGet(string) ([]byte, error)
}

func (odd *OneDriveDescription) Init(api MicrosoftGraphAPI) error {
	bytes, err := api.UseMicrosoftGraphAPIGet("/me/drive")
	if err != nil {
		log.Println(err)
	}
	newMicrosoftGraphDrive := graphapi.MicrosoftGraphDrive{}
	if err := json.Unmarshal(bytes, &newMicrosoftGraphDrive); err != nil {
		log.Println(err)
	}
	if odd.DriveDescription != nil && odd.DriveDescription.ID != newMicrosoftGraphDrive.ID {
		return errors.New("odd.Init DriveDescriptionIDDoesNOTMatchRecord")
	}
	return odd.SetDriveDescription(&newMicrosoftGraphDrive)
}

func (odd *OneDriveDescription) SetDriveDescription(input *graphapi.MicrosoftGraphDrive) error {
	odd.DriveDescription = input
	return nil
}

func (odd *OneDriveDescription) Get() error {
	return nil
}

func (odd *OneDriveDescription) RelativePathToFullDriveRootPath(str string) string {
	return "/me" + odd.RelativePathToDriveRootPath(str)
}

func (odd *OneDriveDescription) RelativePathToDriveRootPath(str string) string {
	path := utils.RegularRootPath(odd.RootPath)
	str = utils.RegularPath(str)
	if str != "/" {
		path += str
	}
	return path
}

func (odd *OneDriveDescription) FullDriveRootPathToRelativePath(str string) string {
	length := len(str)
	if length <= 3 {
		return "/"
	}
	return odd.DriveRootPathToRelativePath(str[3:length])
}

func (odd *OneDriveDescription) DriveRootPathToRelativePath(str string) string {
	rootPath := utils.RegularRootPath(odd.RootPath)
	str, _ = url.QueryUnescape(str)
	strS := strings.Split(str, rootPath)
	path := strS[0]
	for i, s := range strS {
		if i == 1 {
			path += s
		} else if i > 1 {
			path += rootPath + s
		}
	}
	if path == "" {
		path += "/"
	}
	return path
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDrivePath(str string) string {
	return "/me/drive" + str
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDriveItem(str string) string {
	return "/me" + str
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDriveChildren(str string) string {
	if str == "/drive/root:" {
		return "/me/drive/root/children"
	}
	return "/me" + str + ":/children"
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDriveChildrenPath(str string) string {
	return odd.RelativePathToFullDriveRootPath(str) + ":/children"
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDriveExpandChildrenPath(str string) string {
	return odd.RelativePathToFullDriveRootPath(str) + "?expand=children($select=name,size,file,folder,parentReference,createdDateTime,lastModifiedDateTime)"
}

func (odd *OneDriveDescription) UseMicrosoftGraphAPIMeDriveContentPath(str string) string {
	return odd.RelativePathToFullDriveRootPath(str) + ":/content"
}

func (odd *OneDriveDescription) GetRefreshInterval() int64 {
	refreshInterval := int64(2)
	// if odd.RefreshInterval < 60 {
	// 	refreshInterval = int64(60)
	// } else if odd.RefreshInterval < 3000 {
	// 	refreshInterval = odd.RefreshInterval
	// }
	return refreshInterval
}
