package upload

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func (uc *UploaderCollection) Load(driveID string) error {
	uploadCacheFile := driveID + ".upload.json"
	log.Println("Loading OneDrive upload cache file from " + uploadCacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	bytes, err := ioutil.ReadFile(uploadCacheFile)
	if _, ok := err.(*os.PathError); ok {
		log.Println("Creating OneDrive upload cache file " + uploadCacheFile)
		return ioutil.WriteFile(uploadCacheFile, []byte("{}"), 0644)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, uc)
}

func (uc *UploaderCollection) Save(driveID string) error {
	uploadCache := struct {
		Uploaders []*Uploader `json:"uploaders"`
	}{
		uc.Uploaders,
	}

	uploadCacheFile := driveID + ".upload.json"
	bytes, err := json.Marshal(uploadCache)
	if err != nil {
		return err
	}

	log.Println("Saving OneDrive upload cache file to " + uploadCacheFile)
	mutex.Lock()
	defer mutex.Unlock()
	return ioutil.WriteFile(uploadCacheFile, bytes, 0644)
}
