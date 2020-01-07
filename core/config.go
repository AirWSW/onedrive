package core

import (
	"encoding/json"
	"log"
	"os"
)

func CreateOneDriveFromConfigFile() (od OneDrive, err error) {
	argNum := len(os.Args)
	configFile := "config.json"
	if argNum > 1 {
		if os.Args[argNum-2] == "-c" {
			configFile = os.Args[argNum-1]
		}
	}
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&od)
	if err != nil {
		return OneDrive{}, err
	}
	return od, nil
}

func (od *OneDrive) SaveConfigFile() error {
	configFile := "config.json"
	file, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)
	newOD := struct {
		AppRegistrationConfig  AppRegistrationConfig  `json:"appRegistrationConfig"`
		DriveDescriptionConfig DriveDescriptionConfig `json:"driveDescriptionConfig"`
	}{
		od.AppRegistrationConfig,
		od.DriveDescriptionConfig,
	}
	err := encoder.Encode(&newOD)
	if err != nil {
		return err
	}
	log.Println("OneDrive config file saved to " + configFile)
	return nil
}

func (od *OneDrive) HotReloadConfigFile() error {
	argNum := len(os.Args)
	configFile := "config.json"
	if argNum > 1 {
		if os.Args[argNum-2] == "--config" {
			configFile = os.Args[argNum-1]
		}
	}
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)

	newOD := OneDrive{}
	err := decoder.Decode(&newOD)
	if err != nil {
		return err
	}
	od.AppRegistrationConfig = newOD.AppRegistrationConfig
	od.DriveDescriptionConfig = newOD.DriveDescriptionConfig
	log.Println("OneDrive config file hot reloaded")
	return nil
}

func (od *OneDrive) SaveDriveCacheFile() error {
	cacheFile := "cache.json"
	file, _ := os.OpenFile(cacheFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)
	newOD := struct {
		DriveItemsCaches     []DriveItemsCache      `json:"driveItemsCaches"`
		DriveCache           []DriveCache           `json:"driveCache"`
		DriveCacheContentURL []DriveCacheContentURL `json:"driveCacheContentUrl"`
	}{
		od.DriveItemsCaches,
		od.DriveCache,
		od.DriveCacheContentURL,
	}
	err := encoder.Encode(&newOD)
	if err != nil {
		return err
	}
	log.Println("OneDrive cache file saved to " + cacheFile)
	return nil
}

func (od *OneDrive) LoadDriveCacheFile() error {
	cacheFile := "cache.json"
	file, _ := os.Open(cacheFile)
	defer file.Close()
	decoder := json.NewDecoder(file)

	newOD := OneDrive{}
	err := decoder.Decode(&newOD)
	if err != nil {
		return err
	}
	od.DriveItemsCaches = newOD.DriveItemsCaches
	od.DriveCache = newOD.DriveCache
	od.DriveCacheContentURL = newOD.DriveCacheContentURL

	log.Println("OneDrive cache file loaded")
	return nil
}

// func CreateDefaultOneDrive() (od OneDrive, err error) {
// 	od.RegisterConfig = RegisterConfig{
// 		EndPointURI:  "https://login.chinacloudapi.cn/common/oauth2/v2.0/token", // "https://login.microsoftonline.com/common/oauth2/v2.0/token"
// 		Scope:        "user.read files.readwrite.all offline_access",
// 		ClientID:     "",
// 		ClientSecret: "",
// 		RedirectURI:  "http://localhost:8081/auth-redirect",
// 	}
// 	od.DriveConfig = DriveConfig{
// 		EndPointURI:           "https://microsoftgraph.chinacloudapi.cn/v1.0/me/drive", //"https://graph.microsoft.com/v1.0/me/drive"
// 		RootPath:              "root",
// 		FileRefreshInterval:   1200,
// 		FolderRefreshInterval: 600,
// 		Code:                  "",
// 		RefreshToken:          "",
// 	}
// 	return od, nil
// }
