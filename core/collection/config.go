package collection

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/AirWSW/onedrive/core/description"
	"github.com/AirWSW/onedrive/graphapi"
)

var mutex sync.Mutex

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
	return NewOneDriveCollectionFromConfigFile(&ODCollection)
}

func NewOneDriveCollectionFromConfigFile(odc *OneDriveCollection) error {
	return odc.LoadConfigFile()
}

func (odc *OneDriveCollection) LoadConfigFile() error {
	configFile := GetConfigFilenameFromArgs()
	log.Println("Loading OneDriveCollection config file from " + configFile)
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			SaveConfigTemplateFile()
		}
		return err
	}
	return json.Unmarshal(bytes, odc)
}

func (odc *OneDriveCollection) SaveConfigFile() error {
	var newODs []interface{} = nil
	for _, oneDrive := range odc.OneDrives {
		newODs = append(newODs, struct {
			MicrosoftEndPoints     graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
			AzureADAppRegistration graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
			AzureADAuthFlowContext graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
			OneDriveDescription    description.OneDriveDescription `json:"oneDriveDescription"`
		}{
			oneDrive.MicrosoftEndPoints,
			oneDrive.AzureADAppRegistration,
			oneDrive.AzureADAuthFlowContext,
			oneDrive.OneDriveDescription,
		})
	}
	newODC := struct {
		IsDebugMode *bool         `json:"isDebugMode"`
		OneDrives   []interface{} `json:"oneDrives"`
	}{
		odc.IsDebugMode,
		newODs,
	}

	configFile := GetConfigFilenameFromArgs()
	bytes, err := json.MarshalIndent(newODC, "", "    ")
	if err != nil {
		return err
	}

	log.Println("Saving OneDriveCollection config file to " + configFile)
	mutex.Lock()
	defer mutex.Unlock()
	return ioutil.WriteFile(configFile, bytes, 0644)
}

func SaveConfigTemplateFile() error {
	var newODs []interface{} = nil
	newODs = append(newODs, struct {
		MicrosoftEndPoints     graphapi.MicrosoftEndPoints     `json:"microsoftEndPoints"`
		AzureADAppRegistration graphapi.AzureADAppRegistration `json:"azureAdAppRegistration"`
		AzureADAuthFlowContext graphapi.AzureADAuthFlowContext `json:"azureAdAuthFlowContext"`
		OneDriveDescription    description.OneDriveDescription `json:"oneDriveDescription"`
	}{
		graphapi.MicrosoftEndPoints{},
		graphapi.AzureADAppRegistration{
			RedirectURIs: []string{},
		},
		graphapi.AzureADAuthFlowContext{},
		description.OneDriveDescription{
			RefreshInterval: 3600,
		},
	})
	newODC := struct {
		OneDrives []interface{} `json:"oneDrives"`
	}{
		newODs,
	}

	configFile := GetConfigFilenameFromArgs()
	bytes, err := json.MarshalIndent(newODC, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, bytes, 0644); err != nil {
		return err
	}

	log.Println("Creating OneDriveCollection config file template " + configFile)
	mutex.Lock()
	defer mutex.Unlock()
	return ioutil.WriteFile(configFile, bytes, 0644)
}
