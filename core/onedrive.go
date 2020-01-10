package core

import (
	"encoding/json"
	"log"

	"github.com/AirWSW/onedrive/graphapi"
)

type oneDriveCollection interface { // import cycle
	SaveConfigFile() error
}

func (od *OneDrive) Start(odc oneDriveCollection) error { // import cycle
	if err := od.InitMicrosoftGraphAPI(); err != nil {
		return err
	}
	if err := od.InitMicrosoftGraphAPIToken(odc); err != nil {
		if err := odc.SaveConfigFile(); err != nil {
			return err
		}
		return nil
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

func (od *OneDrive) ReStart(odc oneDriveCollection) error { // import cycle
	if err := od.Start(odc); err != nil {
		return err
	}
	if err := odc.SaveConfigFile(); err != nil {
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
	// od.AzureADAuthFlowContext.Code = nil
	return nil
}

func (od *OneDrive) InitMicrosoftGraphAPIToken(odc oneDriveCollection) error { // import cycle
	if od.AzureADAuthFlowContext.RefreshToken == nil {
		if od.AzureADAuthFlowContext.Code == nil {
			if err := od.MicrosoftGraphAPI.GetMicrosoftGraphAPIToken(); err != nil {
				od.AzureADAuthFlowContext.StateID = od.MicrosoftGraphAPI.AzureADAuthFlowContext.StateID
				return err
			}
		}
		if err := od.InitMicrosoftGraphAPI(); err != nil {
			return err
		}
		if err := od.MicrosoftGraphAPI.GetMicrosoftGraphAPIToken(); err != nil {
			return err
		}
		od.AzureADAuthFlowContext.Code = nil
		od.AzureADAuthFlowContext.RefreshToken = od.MicrosoftGraphAPI.AzureADAuthFlowContext.RefreshToken
		if err := odc.SaveConfigFile(); err != nil {
			return err
		}
	} else {
		if err := od.MicrosoftGraphAPI.GetMicrosoftGraphAPIToken(); err != nil {
			return err
		}
	}
	return nil
}

func (od *OneDrive) InitOneDriveDescription() error {
	bytes, err := od.MicrosoftGraphAPI.UseMicrosoftGraphAPIGet("/me/drive")
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
