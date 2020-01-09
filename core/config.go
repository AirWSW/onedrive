package core

import (
	"encoding/json"
	"log"
	"os"

	"github.com/AirWSW/onedrive/graphapi"
)

func GetConfigFilenameFromArgs() string {
	argNum := len(os.Args)
	if argNum > 1 {
		if os.Args[argNum-2] == "-c" {
			return os.Args[argNum-1]
		}
	}
	return "config.json"
}

func InitOneDriveCollectionFromConfigFile() error {
	newODCollection, err := NewOneDriveCollectionFromConfigFile()
	if err != nil {
		return err
	}
	ODCollection = *newODCollection
	return nil
}

func NewOneDriveCollectionFromConfigFile() (*OneDriveCollection, error) {
	configFile := GetConfigFilenameFromArgs()
	file, _ := os.Open(configFile)
	defer file.Close()
	decoder := json.NewDecoder(file)

	log.Println("Loading OneDriveCollection config file from " + configFile)
	odc := &OneDriveCollection{}
	err := decoder.Decode(odc)
	if err != nil {
		return nil, err
	}
	return odc, nil
}

func (odc *OneDriveCollection) SaveConfigFile() error {
	configFile := GetConfigFilenameFromArgs()
	file, _ := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer file.Close()
	encoder := json.NewEncoder(file)

	var newODs []interface{} = nil
	for _, oneDrive := range odc.OneDrives {
		newODs = append(newODs, struct {
			MicrosoftEndPoints     graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
			AzureADAppRegistration graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
			AzureADAuthFlowContext graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
		}{
			oneDrive.MicrosoftEndPoints,
			oneDrive.AzureADAppRegistration,
			oneDrive.AzureADAuthFlowContext,
		})
	}

	newODC := struct {
		OneDrives []interface{} `json:"oneDrives"`
	}{
		newODs,
	}

	log.Println("Saving OneDriveCollection config file to " + configFile)
	return encoder.Encode(&newODC)
}
