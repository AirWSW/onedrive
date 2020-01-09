package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AirWSW/onedrive/core"
)

var OD core.OldOneDrive
var ODCollection *core.OneDriveCollection = &core.ODCollection

func AddDefalutHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
}

func handleGetAuth(c *gin.Context) {
	code := c.Query("code")
	if len(code) > 0 {
		OD.DriveDescriptionConfig.Code = code
		if err := OD.GetMicrosoftGraphAPIToken(); err != nil {
			log.Println(err)
		}
		OD.DriveDescriptionConfig.Code = ""
		if err := OD.SaveConfigFile(); err != nil {
			log.Println(err)
		}
	}
	c.String(http.StatusOK, "code %s", code)
}

func handleGetMicrosoftGraphDriveItem(c *gin.Context) {
	path := c.Query("path")
	microsoftGraphDriveItemCache, err := ODCollection.UseDefaultOneDrive().GetMicrosoftGraphDriveItem(path)
	if err != nil {
		log.Println(err)
	}
	bytes, err := json.Marshal(microsoftGraphDriveItemCache)
	if err != nil {
		log.Println(err)
	}
	AddDefalutHeaders(c)
	c.String(http.StatusOK, "%s", bytes)
}

func handleGetRaw(c *gin.Context) {
	path := c.Query("path")
	log.Println(path)
	bytes, err := ODCollection.OneDrives[0].GetMicrosoftGraphAPIMeDriveRaw(path)
	if err != nil {
		log.Println(err)
	}
	AddDefalutHeaders(c)
	c.String(http.StatusOK, "%s", bytes)
}

func handleGetMicrosoftGraphDriveItemContentURL(c *gin.Context) {
	path := c.Query("path")
	url, err := OD.GetDriveItemContentURLFromPath(path)
	if err != nil {
		log.Println(err)
	}
	AddDefalutHeaders(c)
	c.Redirect(http.StatusFound, url.String())
}

func main() {
	if err := core.InitOneDriveCollectionFromConfigFile(); err != nil {
		log.Panicln(err)
	}
	if err := ODCollection.StartAll(); err != nil {
		log.Panicln(err)
	}
	// if err := core.ODCollection.Run(); err != nil {
	// 	log.Panicln(err)
	// }
	// OD, _ = core.CreateOneDriveFromConfigFile()
	// if err := OD.Run(); err != nil {
	// 	log.Panicln(err)
	// }
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/onedrive/auth", handleGetAuth)
	router.GET("/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/onedrive/raw", handleGetRaw)
	router.GET("/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/auth", handleGetAuth)
	router.GET("/api/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/api/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	if err := router.Run("localhost:8081"); err != nil {
		log.Panicln(err)
	}
}
