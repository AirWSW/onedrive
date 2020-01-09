package core

import (
	"encoding/json"
	"log"

	"github.com/AirWSW/onedrive/graphapi"
)

var ODCollection OneDriveCollection

func (odc *OneDriveCollection) StartAll() error {
	for _, oneDrive := range odc.OneDrives {
		if err := oneDrive.Start(); err != nil {
			return err
		}
	}
	if err := odc.SaveConfigFile(); err != nil {
		return err
	}
	if err := odc.CronStart(); err != nil {
		return err
	}
	return nil
}

func (odc *OneDriveCollection) UseDefaultOneDrive() *OneDrive {
	return odc.OneDrives[0]
}

func (odc *OneDriveCollection) UseOneDrive(str string) *OneDrive {
	for _, oneDrive := range odc.OneDrives {
		if oneDrive.OneDriveDescription.DriveDescription.ID == str {
			return oneDrive
		}
	}
	return nil
}

func (od *OneDrive) Start() error {
	if err := od.InitMicrosoftGraphAPI(); err != nil {
		return err
	}
	if err := od.MicrosoftGraphAPI.GetMicrosoftGraphAPIToken(); err != nil {
		return err
	}
	if err := od.InitOneDriveDescription(); err != nil {
		return err
	}
	if err := od.LoadDriveCacheFile(); err != nil {
		return err
	}
	if err := od.SaveDriveCacheFile(); err != nil {
		return err
	}
	return nil
}

func (od *OneDrive) InitMicrosoftGraphAPI() error {
	input := &graphapi.NewMicrosoftGraphAPIInput{
		MicrosoftEndPoints:     &od.MicrosoftEndPoints,
		AzureADAppRegistration: &od.AzureADAppRegistration,
		AzureADAuthFlowContext: &od.AzureADAuthFlowContext,
	}
	newMicrosoftGraphAPI, err := graphapi.NewMicrosoftGraphAPI(input)
	if err != nil {
		return err
	}
	od.MicrosoftGraphAPI = newMicrosoftGraphAPI
	return nil
}

func (od *OneDrive) InitOneDriveDescription() error {
	bytes, err := od.MicrosoftGraphAPI.UseMicrosoftGraphAPI("/me/drive")
	if err != nil {
		log.Println(err)
	}
	microsoftGraphDrive := graphapi.MicrosoftGraphDrive{}
	if err := json.Unmarshal(bytes, &microsoftGraphDrive); err != nil {
		log.Println(err)
	}
	od.OneDriveDescription.SetDriveDescription(&microsoftGraphDrive)
	return nil
}
