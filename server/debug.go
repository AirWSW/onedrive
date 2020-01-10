package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleGetMicrosoftGraphAPIMeDriveRaw(c *gin.Context) {
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
	bytes, err := od.GetMicrosoftGraphAPIMeDriveRaw(path)
	if err != nil {
		log.Println(err)
		// c.AbortWithStatus(http.StatusNotFound)
		// return
	}
	c.Header("Cache-Control", "private")
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, "%s", bytes)
}

func handlePostMicrosoftGraphAPIMeDriveRaw(c *gin.Context) {
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
	var postBody io.Reader = nil
	var err error = nil
	if c.Request.ContentLength > 0 {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		postBody = bytes.NewReader(data)
	}
	bytes, err := od.PostMicrosoftGraphAPIMeDriveRaw(path, postBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "private")
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, "%s", bytes)
}

func handlePutMicrosoftGraphAPIMeDriveRaw(c *gin.Context) {
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
	bytes, err := od.GetMicrosoftGraphAPIMeDriveRaw(path)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "private")
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, "%s", bytes)
}