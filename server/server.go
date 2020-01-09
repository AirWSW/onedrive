package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AirWSW/onedrive/core"
	"github.com/AirWSW/onedrive/graphapi"
)

var OD core.OldOneDrive
var ODCollection *core.OneDriveCollection = &core.ODCollection

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

func handleGetDrive(c *gin.Context) {
	path := c.Query("path")
	driveCache, err := OD.GetDriveItemsFromPath(path)
	if err != nil {
		log.Println(err)
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	bytes, _ := json.Marshal(driveCache)
	c.String(http.StatusOK, "%s", bytes)
}

func handleGetRaw(c *gin.Context) {
	path := c.Query("path")
	log.Println(path)
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Add("Authorization", "Bearer "+OD.MicrosoftGraphAPIToken.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	v := &graphapi.MicrosoftGraphDrive{}
	if err = json.Unmarshal([]byte(body), v); err != nil {
		log.Println(err)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	c.String(http.StatusOK, "%s", body)
}

func handleGetFile(c *gin.Context) {
	path := c.Query("path")
	url, err := OD.GetDriveItemContentURLFromPath(path)
	if err != nil {
		log.Println(err)
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	c.Redirect(http.StatusFound, url.String())
}

func handleGetStream(c *gin.Context) {
	path := c.Param("path")
	url, err := OD.GetDriveItemContentURLFromPath(path)
	if err != nil {
		log.Println(err)
	}
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	c.Redirect(http.StatusFound, url.String())
}

func main() {
	if err := core.InitOneDriveCollectionFromConfigFile(); err != nil {
		log.Panicln(err)
	}
	if err := ODCollection.StartAll(); err != nil {
		log.Panicln(err)
	}

	log.Println(ODCollection)
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
	router.GET("/onedrive/drive", handleGetDrive)
	router.GET("/onedrive/file", handleGetFile)
	router.GET("/onedrive/raw", handleGetRaw)
	router.GET("/onedrive/stream/*path", handleGetStream)
	router.GET("/api/onedrive/auth", handleGetAuth)
	router.GET("/api/onedrive/drive", handleGetDrive)
	router.GET("/api/onedrive/file", handleGetFile)
	router.GET("/api/onedrive/stream/*path", handleGetStream)
	if err := router.Run("localhost:8081"); err != nil {
		log.Panicln(err)
	}
}
