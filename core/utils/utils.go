package utils

import (
	"net/url"
	"strings"
)

func RegularPathToPathFilename(str string) (path, filename string) {
	strS := strings.Split(str, "/")
	strR := ""
	for _, s := range strS {
		if s != "" {
			strR += "/" + s
		}
	}
	if strR == "" {
		strR += "/"
	}
	strRR := strings.Split(strR, "/")
	n := len(strRR)
	if n > 1 {
		for _, s := range strRR[0 : n-1] {
			if s != "" {
				path += "/" + s
			}
		}
		filename = strRR[n-1]
	}
	return path, filename
}

func RegularPath(str string) (path string) {
	// any path to "/" or "/path/to"
	str, _ = url.QueryUnescape(str)
	strD := strings.Split(str, "#")
	strQ := strings.Split(strD[0], "?")
	pathRaw := strQ[0]
	pathQuery := ""
	for i, s := range strQ {
		if i != 0 {
			pathQuery += s
		}
	}
	strS := strings.Split(pathRaw, "/")
	for _, s := range strS {
		if s != "" {
			path += "/" + s
		}
	}
	if path == "" {
		path += "/"
	}
	if pathQuery != "" {
		path += "?" + pathQuery
	}
	return path
}

func RegularRootPath(str string) (path string) {
	length := len(str)
	if length > 0 {
		str = RegularPath(str)
		if str == "/" || str == "/root" || str == "/root:" {
			path = "/drive/root:"
		} else {
			path = "/drive/root:" + str
		}
	} else {
		path = "/drive/root:"
	}
	return path
}
