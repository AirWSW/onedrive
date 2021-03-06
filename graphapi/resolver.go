package graphapi

import uuid "github.com/satori/go.uuid"

func (e *MicrosoftEndPoints) Set(input *MicrosoftEndPoints) error {
	e.AzureADPortalEndPointURL = input.AzureADPortalEndPointURL
	e.AzureADEndPointURL = input.AzureADEndPointURL
	e.MicrosoftGraphAPIEndPointURL = input.MicrosoftGraphAPIEndPointURL
	return nil
}

func (e *MicrosoftEndPoints) GetAzureADAuthorizeEndPointURL() string {
	return e.AzureADEndPointURL + "/common/oauth2/v2.0/authorize"
}

func (e *MicrosoftEndPoints) PostAzureADTokenEndPointURL() string {
	return e.AzureADEndPointURL + "/common/oauth2/v2.0/token"
}

func (e *MicrosoftEndPoints) GetMicrosoftGraphAPIEndPointURL() string {
	return e.MicrosoftGraphAPIEndPointURL + "/v1.0"
}

func (e *MicrosoftEndPoints) UseMicrosoftGraphAPIEndPointURL(str string) string {
	return e.GetMicrosoftGraphAPIEndPointURL() + str
}

func (r *AzureADAppRegistration) Set(input *AzureADAppRegistration) error {
	r.DisplayName = input.DisplayName
	r.ClientID = input.ClientID
	r.TenantID = input.TenantID
	r.ObjectID = input.ObjectID
	r.RedirectURIs = input.RedirectURIs
	r.LogoutURL = input.LogoutURL
	r.ClientSecret = input.ClientSecret
	return nil
}

func (c *AzureADAuthFlowContext) Set(input *AzureADAuthFlowContext) error {
	c.GrantScope = input.GrantScope
	if input.StateID == nil {
		stateID := uuid.Must(uuid.NewV4(), nil).String()
		c.StateID = &stateID
	} else {
		c.StateID = input.StateID
	}
	c.Code = input.Code
	c.RefreshToken = input.RefreshToken
	return nil
}

func (c *AzureADAuthFlowContext) SetRefreshToken(input *string) error {
	c.RefreshToken = input
	return nil
}

func (t *MicrosoftGraphAPIToken) Set(input *MicrosoftGraphAPIToken) error {
	t.TokenType = input.TokenType
	t.ExpiresIn = input.ExpiresIn
	t.ExtExpiresIn = input.ExtExpiresIn
	t.Scope = input.Scope
	t.AccessToken = input.AccessToken
	t.RefreshToken = input.RefreshToken
	return nil
}

func (t *MicrosoftGraphAPIToken) GetAuthorizationString() string {
	return t.TokenType + " " + t.AccessToken
}
