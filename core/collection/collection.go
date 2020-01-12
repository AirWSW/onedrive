package collection

import "github.com/AirWSW/onedrive/core"

var ODCollection OneDriveCollection

func (odc *OneDriveCollection) StartAll() error {
	for _, oneDrive := range odc.OneDrives {
		if err := oneDrive.Start(odc); err != nil {
			return err
		}
	}
	if err := odc.SaveConfigFile(); err != nil {
		return err
	}
	if err := odc.CronStartAll(); err != nil {
		return err
	}
	return nil
}

func (odc *OneDriveCollection) UseDefaultOneDrive() *core.OneDrive {
	return odc.OneDrives[0]
}

func (odc *OneDriveCollection) UseOneDriveByID(str string) *core.OneDrive {
	for _, oneDrive := range odc.OneDrives {
		if oneDrive.OneDriveDescription.DriveDescription != nil && oneDrive.OneDriveDescription.DriveDescription.ID == str {
			return oneDrive
		}
	}
	return nil
}

func (odc *OneDriveCollection) UseOneDriveByOneDriveName(str string) *core.OneDrive {
	for _, oneDrive := range odc.OneDrives {
		if oneDrive.OneDriveDescription.OneDriveName != nil && *oneDrive.OneDriveDescription.OneDriveName == str {
			return oneDrive
		}
	}
	return nil
}

func (odc *OneDriveCollection) UseOneDriveByStateID(str string) *core.OneDrive {
	for _, oneDrive := range odc.OneDrives {
		if oneDrive.AzureADAuthFlowContext.StateID != nil && *oneDrive.AzureADAuthFlowContext.StateID == str {
			return oneDrive
		}
	}
	return nil
}
