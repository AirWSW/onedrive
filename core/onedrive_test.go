package core

import (
	"net/url"
	"strings"
	"testing"
)

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
