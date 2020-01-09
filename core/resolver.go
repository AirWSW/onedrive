package core

import (
	"net/url"
	"strings"

	"github.com/AirWSW/onedrive/graphapi"
)

func (odd *OneDriveDescription) SetDriveDescription(input *graphapi.MicrosoftGraphDrive) error {
	odd.DriveDescription = input
	return nil
}

func (odd *OneDriveDescription) Get() error {
	return nil
}

func RegularPathToPathFilename(str string) (path, filename string) {
	strS := strings.Split(str, "/")
	strR := ""
	for _, s := range strS {
		if s != "" {
			strR += "/" + s
		}
	}
	if strR == "" {
		strR += "/"
	}
	strRR := strings.Split(strR, "/")
	n := len(strRR)
	if n > 1 {
		for _, s := range strRR[0 : n-1] {
			if s != "" {
				path += "/" + s
			}
		}
		filename = strRR[n-1]
	}
	return path, filename
}

func RegularPath(str string) (path string) {
	// any path to "/" or "/path/to"
	str, _ = url.QueryUnescape(str)
	strD := strings.Split(str, "#")
	strQ := strings.Split(strD[0], "?")
	pathRaw := strQ[0]
	pathQuery := ""
	for i, s := range strQ {
		if i != 0 {
			pathQuery += s
		}
	}
	strS := strings.Split(pathRaw, "/")
	for _, s := range strS {
		if s != "" {
			path += "/" + s
		}
	}
	if path == "" {
		path += "/"
	}
	if pathQuery != "" {
		path += "?" + pathQuery
	}
	return path
}

func RegularRootPath(str string) (path string) {
	length := len(str)
	if length > 0 {
		str = RegularPath(str)
		if str == "/" || str == "/root" || str == "/root:" {
			path = "/drive/root:"
		} else {
			path = "/drive/root:" + str
		}
	} else {
		path = "/drive/root:"
	}
	return path
}

func (odd *OneDriveDescription) RelativePathToFullDriveRootPath(str string) string {
	return "/me" + odd.RelativePathToDriveRootPath(str)
}

func (odd *OneDriveDescription) RelativePathToDriveRootPath(str string) string {
	path := RegularRootPath(odd.RootPath)
	str = RegularPath(str)
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
	rootPath := RegularRootPath(odd.RootPath)
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
