package core

import (
	"net/url"
	"strings"
	"testing"
)

// func TestCreateOneDriveCollectionFromConfigFile(t *testing.T) {
// 	odc, err := NewOneDriveCollectionFromConfigFile()
// 	if err != nil {
// 		t.Fatalf("%s", err)
// 	}

// 	ods := odc.OneDrives
// 	od := ods[0]
// 	if od.AzureADAppRegistration.DisplayName != nil {
// 		t.Logf("%s", *od.AzureADAppRegistration.DisplayName)
// 	}
// 	t.Logf("%s", od.AzureADAppRegistration.ClientID)
// 	if od.AzureADAppRegistration.TenantID != nil {
// 		t.Logf("%s", *od.AzureADAppRegistration.TenantID)
// 	}
// 	if od.AzureADAppRegistration.ObjectID != nil {
// 		t.Logf("%s", *od.AzureADAppRegistration.ObjectID)
// 	}
// 	if od.AzureADAppRegistration.ObjectID != nil {
// 		t.Logf("%s", *od.AzureADAppRegistration.ObjectID)
// 	}
// 	t.Logf("%s", od.AzureADAppRegistration.RedirectURIs)
// 	if od.AzureADAppRegistration.LogoutURL != nil {
// 		t.Logf("%s", *od.AzureADAppRegistration.LogoutURL)
// 	}
// 	t.Logf("%s", od.AzureADAppRegistration.ClientSecret)
// 	if od.MicrosoftEndPoints.AzureADPortalEndPointURL != nil {
// 		t.Logf("%s", *od.MicrosoftEndPoints.AzureADPortalEndPointURL)
// 	}
// 	t.Logf("%s", od.MicrosoftEndPoints.AzureADEndPointURL)
// 	t.Logf("%s", od.MicrosoftEndPoints.MicrosoftGraphAPIEndPointURL)

// 	input := &NewMicrosoftGraphAPIInput{
// 		MicrosoftEndPoints:     &od.MicrosoftEndPoints,
// 		AzureADAppRegistration: &od.AzureADAppRegistration,
// 		AzureADAuthFlowContext: &od.AzureADAuthFlowContext,
// 	}
// 	od.MicrosoftGraphAPI, err = NewMicrosoftGraphAPI(input)
// 	if err != nil {
// 		t.Fatalf("%s", err)
// 	}

// 	od.MicrosoftGraphAPI.GetMicrosoftGraphAPIToken()
// 	t.Logf("%s", od.MicrosoftGraphAPI.MicrosoftGraphAPIToken.AccessToken)
// }

func TestRegularPath(t *testing.T) {
	strO := ""
	// strO := "//video.airw.cf/tv.shows////joy.of.life.s01.2019//////e02?t=123456?123&123456"
	t.Logf(strO)
	strD := strings.Split(strO, "#")
	strQ := strings.Split(strD[0], "?")
	pathRaw := strQ[0]
	pathQuery := ""
	for i, str := range strQ {
		if i != 0 {
			pathQuery += str
		}
		t.Logf("%d %s", i, str)
	}
	t.Logf("pathRaw %s", pathRaw)
	t.Logf("pathQuery %s", pathQuery)

	strS := strings.Split(pathRaw, "/")
	path := ""
	for i, str := range strS {
		if str != "" {
			path += "/" + str
		}
		t.Logf("%d %s", i, str)
	}
	if path == "" {
		path += "/"
	}
	if pathQuery != "" {
		path += "?" + pathQuery
	}
	t.Logf("path %s", path)
}

func TestPathToFileName(t *testing.T) {
	strO := "//video.airw.cf/tv.shows////joy.of.life.s01.2019//////e02/0000.ts"
	// strO = ""
	t.Logf(strO)
	strS := strings.Split(strO, "/")
	strR := ""
	for i, str := range strS {
		if str != "" {
			strR += "/" + str
		}
		t.Logf("%d %s", i, str)
	}
	if strR == "" {
		strR += "/"
	}
	t.Logf("strR %s", strR)

	strRR := strings.Split(strR, "/")
	n := len(strRR)
	path := ""
	filename := ""
	if n > 1 {
		for i, str := range strRR[0 : n-1] {
			if str != "" {
				path += "/" + str
			}
			t.Logf("%d %s", i, str)
		}
		filename = strRR[n-1]
	}
	t.Logf("path %s", path)
	t.Logf("filename %s", filename)
}

func TestURL(t *testing.T) {
	strO := "https://login.chinacloudapi.cn/common/oauth2/v2.0/token"
	t.Logf(strO)
	endPointURI, _ := url.Parse(strO)
	t.Logf("path %s", endPointURI.Host)
}

func TestURL2(t *testing.T) {
	strO := "/me"
	length := len(strO)
	if length <= 3 {
		t.Logf("/")
	}
	t.Logf(strO[3:length])
}

func TestURL3(t *testing.T) {
	strO := "od://username:password@od"
	myURL, err := url.Parse(strO)
	if err != nil {
		t.Logf("err")
	}
	t.Logf(myURL.String())
}
