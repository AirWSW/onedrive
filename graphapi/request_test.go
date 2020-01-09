package graphapi

import (
	"net/url"
	"testing"
)

func TestNewMicrosoftGraphAPI(t *testing.T) {
	urlStr := "https://login.microsoftonline.com/common/oauth2/v2.0"
	myUrl, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("%s", err)
	}
	t.Logf("%s://%s", myUrl.Scheme, myUrl.Host)
}
