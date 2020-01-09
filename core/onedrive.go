package core

import (
	"github.com/AirWSW/onedrive/graphapi"
)

var ODCollection OneDriveCollection

func (odc *OneDriveCollection) StartAll() error {
	for _, oneDrive := range odc.OneDrives {
		if err := oneDrive.Start(); err != nil {
			return err
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
	// if err := od.SaveConfigFile(); err != nil {
	// 	return err
	// }
	// if err := od.LoadDriveCacheFile(); err != nil {
	// 	return err
	// }
	// if err := od.Cron(); err != nil {
	// 	return err
	// }
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

func (odd *OneDriveDescription) GetDriveDescription() error {

	return nil
}
