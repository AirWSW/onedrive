package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AirWSW/onedrive"

	"github.com/gin-gonic/gin"
)

var OD onedrive.OneDrive

func handleGetAuthRedirect(c *gin.Context) {
	code := c.Query("code")
	c.String(http.StatusOK, "code %s", code)
}

func handleGetAPIOneDrive(c *gin.Context) {
	path := c.Query("path")
	driveCache, err := OD.GetDrivePath(path)
	log.Println(err)
	bytes, _ := json.Marshal(driveCache)
	c.String(http.StatusOK, "%s", bytes)
}

func handleGetAPIDownload(c *gin.Context) {
	path := c.Query("path")
	url, err := OD.GetDrivePathContentURL(path)
	log.Println(err)
	c.Redirect(http.StatusFound, url.String())
}

func main() {
	OD, _ = onedrive.CreateDefaultOneDrive()
	if err := OD.GetMSGraphAPIToken(); err != nil {
		log.Panicln(err)
	}
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/auth-redirect", handleGetAuthRedirect)
	router.GET("/api/onedrive", handleGetAPIOneDrive)
	router.GET("/api/download", handleGetAPIDownload)
	if err := router.Run(":8081"); err != nil {
		log.Panicln(err)
	}
}
