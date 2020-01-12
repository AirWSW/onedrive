package main

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"

	"github.com/AirWSW/onedrive/core/collection"
)

var ODCollection *collection.OneDriveCollection = &collection.ODCollection

func AddDefalutHeaders(c *gin.Context) {
	// c.Header("Access-Control-Allow-Origin", "*")
	// c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	// c.Header("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
	// c.Header("Access-Control-Expose-Headers", "Content-Length,Content-Range")
	c.Header("Cache-Control", "private")
	c.Header("Content-Type", "application/json;charset=utf-8")
}

func handleGetOneDriveStatus(c *gin.Context) {
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		AddDefalutHeaders(c)
		c.String(http.StatusOK, "%s", []byte(`{"status":"ok","drive":"`+drive+`"}`))
		return
	}
	AddDefalutHeaders(c)
	c.String(http.StatusOK, "%s", []byte(`{"status":"ok"}`))
}

func handleGetAzureADAuth(c *gin.Context) {
	state := c.Query("state")
	if len(state) > 0 {
		od := ODCollection.UseOneDriveByStateID(state)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		code := c.Query("code")
		if len(code) > 0 && od.AzureADAuthFlowContext.RefreshToken == nil {
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
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	path := c.Query("path")
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
		AddDefalutHeaders(c)
		c.String(http.StatusOK, "%s", bytes)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func handleGetMicrosoftGraphDriveItemSearch(c *gin.Context) {
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	query := c.Query("query")
	microsoftGraphDriveItemCache, err := od.GetMicrosoftGraphDriveItem(query)
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
		AddDefalutHeaders(c)
		c.String(http.StatusOK, "%s", bytes)
		return
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func handleGetMicrosoftGraphDriveItemContentURL(c *gin.Context) {
	drive := c.Query("drive")
	od := ODCollection.UseDefaultOneDrive()
	if len(drive) > 0 {
		od = ODCollection.UseOneDriveByOneDriveName(drive)
		if od == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
	path := c.Query("path")
	if path == "" {
		path = c.Param("path")
	}
	microsoftGraphDriveItemCache, err := od.GetMicrosoftGraphAPIMeDriveContentURL(path)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if microsoftGraphDriveItemCache != nil {
		if microsoftGraphDriveItemCache.DownloadURL != nil {
			AddDefalutHeaders(c)
			c.Redirect(http.StatusFound, *microsoftGraphDriveItemCache.DownloadURL)
			return
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func main() {
	if ODCollection.IsDebugMode != nil && *ODCollection.IsDebugMode {
		runtime.GOMAXPROCS(1)
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}
	if err := collection.InitOneDriveCollectionFromConfigFile(); err != nil {
		log.Panicln(err)
	}
	if err := ODCollection.StartAll(); err != nil {
		log.Panicln(err)
	}
	router := gin.Default()
	if ODCollection.IsDebugMode != nil && *ODCollection.IsDebugMode {
		gin.SetMode(gin.DebugMode)
		ginpprof.Wrap(router)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if ODCollection.IsDebugMode != nil && *ODCollection.IsDebugMode {
		router.GET("/onedrive/raw", handleGetMicrosoftGraphAPIMeDriveRaw)
		router.POST("/onedrive/raw", handlePostMicrosoftGraphAPIMeDriveRaw)
		router.PUT("/onedrive/raw", handlePutMicrosoftGraphAPIMeDriveRaw)
	}
	router.GET("/onedrive/auth", handleGetAzureADAuth)
	router.GET("/onedrive/content", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/onedrive/driveitem", handleGetMicrosoftGraphDriveItem)
	router.GET("/onedrive/search", handleGetMicrosoftGraphDriveItemSearch)
	router.GET("/onedrive/status", handleGetOneDriveStatus)
	router.GET("/api/onedrive/auth", handleGetAzureADAuth)
	router.GET("/api/onedrive/content", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/driveitem", handleGetMicrosoftGraphDriveItem)
	router.GET("/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/drive", handleGetMicrosoftGraphDriveItem)
	router.GET("/api/onedrive/file", handleGetMicrosoftGraphDriveItemContentURL)
	router.GET("/api/onedrive/stream/*path", handleGetMicrosoftGraphDriveItemContentURL)
	if err := router.Run("localhost:8081"); err != nil {
		log.Panicln(err)
	}
}
