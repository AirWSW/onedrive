package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AirWSW/onedrive/core/collection"
)

var ODCollection *collection.OneDriveCollection = &collection.ODCollection

func AddDefalutHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	c.Next()
}

func handleGetAuth(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if len(state) > 0 {
		od := ODCollection.UseOneDriveByStateID(state)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if len(code) > 0 {
			od.AzureADAuthFlowContext.Code = &code
			if err := od.ReStart(ODCollection); err != nil {
				log.Println(err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.String(http.StatusOK, "code %s", code)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func handleGetMicrosoftGraphDriveItem(c *gin.Context) {
	path := c.Query("path")
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	microsoftGraphDriveItemCache, err := od.GetMicrosoftGraphDriveItem(path)
	if err != nil {
		log.Println(err)
	}
	if microsoftGraphDriveItemCache != nil {
		bytes, err := json.Marshal(microsoftGraphDriveItemCache)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
		c.String(http.StatusOK, "%s", bytes)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func handleGetMicrosoftGraphDriveItemContentURL(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = c.Param("path")
	}
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	microsoftGraphDriveItemCache, err := od.GetMicrosoftGraphAPIMeDriveContentURL(path)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if microsoftGraphDriveItemCache != nil {
		if microsoftGraphDriveItemCache.DownloadURL != nil {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
			c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
			c.Redirect(http.StatusFound, *microsoftGraphDriveItemCache.DownloadURL)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func main() {
	if err := collection.InitOneDriveCollectionFromConfigFile(); err != nil {
		log.Panicln(err)
	}
	if err := ODCollection.StartAll(); err != nil {
		log.Panicln(err)
	}
	if ODCollection.IsDebugMode != nil && *ODCollection.IsDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	if ODCollection.IsDebugMode != nil && *ODCollection.IsDebugMode {
		router.GET("/onedrive/raw", handleGetMicrosoftGraphAPIMeDriveRaw)
		router.POST("/onedrive/raw", handlePostMicrosoftGraphAPIMeDriveRaw)
		router.PUT("/onedrive/raw", handlePutMicrosoftGraphAPIMeDriveRaw)
	}
	router.GET("/onedrive/auth", handleGetAuth)
	router.GET("/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/auth", handleGetAuth)
	router.GET("/api/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/api/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	if err := router.Run("localhost:8081"); err != nil {
		log.Panicln(err)
	}
}
