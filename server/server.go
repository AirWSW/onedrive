package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AirWSW/onedrive/core"
)

var OD core.OneDrive

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
	driveCache, err := OD.GetDrivePath(path)
	if err != nil {
		log.Println(err)
	}
	bytes, _ := json.Marshal(driveCache)
	c.String(http.StatusOK, "%s", bytes)
}

func handleGetFile(c *gin.Context) {
	path := c.Query("path")
	url, err := OD.GetDrivePathContentURL(path)
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
	OD, _ = core.CreateOneDriveFromConfigFile()
	if err := OD.Run(); err != nil {
		log.Panicln(err)
	}
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/auth", handleGetAuth)
	router.GET("/drive", handleGetDrive)
	router.GET("/file", handleGetFile)
	router.GET("/onedrive/auth", handleGetAuth)
	router.GET("/onedrive/drive", handleGetDrive)
	router.GET("/onedrive/file", handleGetFile)
	if err := router.Run("localhost:8081"); err != nil {
		log.Panicln(err)
	}
}
