package graphapi

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NewMicrosoftGraphAPI validates NewMicrosoftGraphAPIInput and assigns to api
func NewMicrosoftGraphAPI(input *NewMicrosoftGraphAPIInput) (*MicrosoftGraphAPI, error) {
	api := &MicrosoftGraphAPI{
		MicrosoftGraphAPIToken: &MicrosoftGraphAPIToken{},
	}

	// Validation input MicrosoftEndPoints and assign to api
	var newAzureADPortalEndPointURL *string = nil
	microsoftEndPoints := input.MicrosoftEndPoints
	if microsoftEndPoints.AzureADPortalEndPointURL != nil {
		myURL, err := url.Parse(*microsoftEndPoints.AzureADPortalEndPointURL)
		if err != nil {
			return nil, err
		}
		if myURL.Scheme != "https" || myURL.Host == "" {
			return nil, errors.New("Invalid AzureADPortalEndPointURL input")
		}
		urlString := myURL.Scheme + "://" + myURL.Host
		newAzureADPortalEndPointURL = &urlString
	}
	var newAzureADEndPointURL, newMicrosoftGraphAPIEndPointURL string = "", ""
	if microsoftEndPoints.AzureADEndPointURL != "" {
		myURL, err := url.Parse(microsoftEndPoints.AzureADEndPointURL)
		if err != nil {
			return nil, err
		}
		if myURL.Scheme != "https" || myURL.Host == "" {
			return nil, errors.New("Invalid AzureADEndPointURL input")
		}
		urlString := myURL.Scheme + "://" + myURL.Host
		newAzureADEndPointURL = urlString
	} else {
		return nil, errors.New("Must input AzureADEndPointURL")
	}
	if microsoftEndPoints.MicrosoftGraphAPIEndPointURL != "" {
		myURL, err := url.Parse(microsoftEndPoints.MicrosoftGraphAPIEndPointURL)
		if err != nil {
			return nil, err
		}
		if myURL.Scheme != "https" || myURL.Host == "" {
			return nil, errors.New("Invalid MicrosoftGraphAPIEndPointURL input")
		}
		urlString := myURL.Scheme + "://" + myURL.Host
		newMicrosoftGraphAPIEndPointURL = urlString
	} else {
		return nil, errors.New("Must input MicrosoftGraphAPIEndPointURL")
	}
	if err := api.MicrosoftEndPoints.Set(&MicrosoftEndPoints{
		AzureADPortalEndPointURL:     newAzureADPortalEndPointURL,
		AzureADEndPointURL:           newAzureADEndPointURL,
		MicrosoftGraphAPIEndPointURL: newMicrosoftGraphAPIEndPointURL,
	}); err != nil {
		return nil, err
	}

	// Validation input AzureADAppRegistration and assign to api
	azureADAppRegistration := input.AzureADAppRegistration
	if azureADAppRegistration.ClientID == "" || azureADAppRegistration.ClientSecret == "" || azureADAppRegistration.RedirectURIs == nil {
		return nil, errors.New("Invalid AzureADAppRegistration input")
	}
	if len(azureADAppRegistration.RedirectURIs) == 0 {
		return nil, errors.New("Must input at least one RedirectURI")
	}
	if err := api.AzureADAppRegistration.Set(&AzureADAppRegistration{
		DisplayName:  azureADAppRegistration.DisplayName,
		ClientID:     azureADAppRegistration.ClientID,
		TenantID:     azureADAppRegistration.TenantID,
		ObjectID:     azureADAppRegistration.ObjectID,
		RedirectURIs: azureADAppRegistration.RedirectURIs,
		LogoutURL:    azureADAppRegistration.LogoutURL,
		ClientSecret: azureADAppRegistration.ClientSecret,
	}); err != nil {
		return nil, err
	}

	// Validation input AzureADAuthFlowContext and assign to api
	azureADAuthFlowContext := input.AzureADAuthFlowContext
	if azureADAuthFlowContext.GrantScope == "" {
		return nil, errors.New("Must input GrantScope")
	}
	if err := api.AzureADAuthFlowContext.Set(&AzureADAuthFlowContext{
		GrantScope:   azureADAuthFlowContext.GrantScope,
		Code:         azureADAuthFlowContext.Code,
		RefreshToken: azureADAuthFlowContext.RefreshToken,
	}); err != nil {
		return nil, err
	}

	// return *MicrosoftGraphAPI as api
	return api, nil
}

func (api *MicrosoftGraphAPI) getMicrosoftGraphAPITokenRequestPostForm() (io.Reader, error) {
	data := url.Values{}
	azureADAuthFlowContext := api.AzureADAuthFlowContext
	azureADAppRegistration := api.AzureADAppRegistration

	// Try RefreshToken and Code
	if azureADAuthFlowContext.RefreshToken != nil {
		data.Set("grant_type", "refresh_token")
		data.Set("refresh_token", *azureADAuthFlowContext.RefreshToken)
	} else if azureADAuthFlowContext.Code != nil {
		data.Set("grant_type", "authorization_code")
		data.Set("code", *azureADAuthFlowContext.Code)
	} else {
		// If both RefreshToken and Code are invalid, log error and return authorize urls
		log.Println("Invalid Microsoft Graph API Token Grant Type, use the following URLs to GET code")
		clientID := azureADAppRegistration.ClientID
		grantScope := url.QueryEscape(azureADAuthFlowContext.GrantScope)
		getAzureADAuthorizeEndPointURL := api.MicrosoftEndPoints.GetAzureADAuthorizeEndPointURL()
		for _, redirectURI := range azureADAppRegistration.RedirectURIs {
			log.Println(getAzureADAuthorizeEndPointURL + "?client_id=" + clientID + "&scope=" + grantScope + "&response_type=code&redirect_uri=" + redirectURI)
		}
		return nil, errors.New("Invalid Microsoft Graph API Token Grant Type")
	}

	// Setting other post form data
	data.Set("client_id", azureADAppRegistration.ClientID)
	data.Set("client_secret", azureADAppRegistration.ClientSecret)
	redirectURIs := azureADAppRegistration.RedirectURIs
	data.Set("redirect_uri", redirectURIs[0])

	// return io.Reader
	return strings.NewReader(data.Encode()), nil
}

func (api *MicrosoftGraphAPI) getMicrosoftGraphAPITokenRequest() error {
	// New post request
	postAzureADTokenEndPointURL := api.MicrosoftEndPoints.PostAzureADTokenEndPointURL()
	postForm, err := api.getMicrosoftGraphAPITokenRequestPostForm()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", postAzureADTokenEndPointURL, postForm)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Get post response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// json.Unmarshal the post response
	newMicrosoftGraphAPIToken := &MicrosoftGraphAPIToken{}
	if err = json.Unmarshal(body, newMicrosoftGraphAPIToken); err != nil {
		return err
	}
	if newMicrosoftGraphAPIToken == nil {
		log.Println(string(body))
		return errors.New("GetMicrosoftGraphAPITokenRequestError")
	}
	if err := api.MicrosoftGraphAPIToken.Set(newMicrosoftGraphAPIToken); err != nil {
		return err
	}

	// Bind api.MicrosoftGraphAPIToken.RefreshToken to api.AzureADAuthFlowContext.RefreshToken
	return api.AzureADAuthFlowContext.SetRefreshToken(api.MicrosoftGraphAPIToken.RefreshToken)
}

func (api *MicrosoftGraphAPI) GetMicrosoftGraphAPIToken() error {
	return api.getMicrosoftGraphAPITokenRequest()
}

func (api *MicrosoftGraphAPI) RefreshMicrosoftGraphAPIToken() error {
	return api.GetMicrosoftGraphAPIToken()
}