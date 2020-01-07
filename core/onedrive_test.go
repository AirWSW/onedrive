package core

import (
	"strings"
	"testing"
)

func TestRegularPath(t *testing.T) {
	strO := "2322/g/g//g/g//123/"
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
