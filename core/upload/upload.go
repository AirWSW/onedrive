package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/AirWSW/onedrive/graphapi"
)

var mutex sync.Mutex

type MicrosoftGraphAPI interface {
	UseMicrosoftGraphAPIGet(string) ([]byte, error)
	UseMicrosoftGraphAPIPost(string, io.Reader) ([]byte, error)
	UseMicrosoftGraphAPIPut(string, io.Reader) ([]byte, error)
}

func (uc *UploaderCollection) Init(api MicrosoftGraphAPI) error {
	filename := "04"
	uploaderReference := &UploaderReference{
		DriveType: "business",
		Size:      629145600,
		Name:      filename,
		Path:      "/drive/root:/cdn/post/uploder/test/" + filename,
	}
	atMicrosoftGraphConflictBehavior := "replace"
	uploadableProperties := &graphapi.MicrosoftGraphDriveItemUploadableProperties{
		Name:                             filename,
		AtMicrosoftGraphConflictBehavior: &atMicrosoftGraphConflictBehavior,
	}
	input := &UploaderDescription{
		UploaderReference:    uploaderReference,
		UploadableProperties: uploadableProperties,
	}
	uploader, err := NewUploader(input)
	if err != nil {
		return err
	}
	uploader.Start(api)
	defer uploader.Close(api)

	return nil
}

func NewUploader(input *UploaderDescription) (*Uploader, error) {
	uploaderDescription := &UploaderDescription{
		UploaderReference:    input.UploaderReference,
		UploadableProperties: input.UploadableProperties,
	}
	uploader := &Uploader{
		UploaderDescription: uploaderDescription,
	}
	return uploader, nil
}

func (u *Uploader) Start(api MicrosoftGraphAPI) {
	uploaderDescription := u.UploaderDescription
	path := UseMicrosoftGraphAPIMeDrivecreateUploadSessionPath(uploaderDescription.UploaderReference.Path)
	data, err := json.Marshal(uploaderDescription.UploadableProperties)
	if err != nil {
		log.Println(err)
	}
	postBody := bytes.NewReader(data)
	respBody, err := api.UseMicrosoftGraphAPIPost(path, postBody)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(respBody))
	microsoftGraphUploadSession := &graphapi.MicrosoftGraphUploadSession{}
	if err := json.Unmarshal(respBody, microsoftGraphUploadSession); err != nil {
		log.Println(err)
	}
	uploadURL := microsoftGraphUploadSession.UploadURL
	size := u.UploaderDescription.UploaderReference.Size
	u.UploaderDescription.UploaderReference.UploadURL = uploadURL
	uploadSession, err := NewUploadSession(size, microsoftGraphUploadSession)
	if err != nil {
		log.Println(err)
	}
	payload := strings.NewReader(randSeq(uploadSession.UploadSessionDescription.GetContentChunkSizeInt64()))
	microsoftGraphUploadSession, err = uploadSession.Put(*uploadURL, payload)
	if err != nil {
		log.Println(err)
	}
	for {
		uploadSessions, err := NewUploadSessionsFromRange(size, microsoftGraphUploadSession)
		if err != nil {
			log.Println(err)
		}
		for _, innerUploadSession := range uploadSessions {
			payload = strings.NewReader(randSeq(innerUploadSession.UploadSessionDescription.GetContentChunkSizeInt64()))
			microsoftGraphUploadSession, err = innerUploadSession.Put(*uploadURL, payload)
			if err != nil {
				log.Println(err)
			}
		}
	}
	log.Println("v")

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int64) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	// log.Println(string(b))
	return string(b)
}

func NewUploadSession(size int64, microsoftGraphUploadSession *graphapi.MicrosoftGraphUploadSession) (*UploadSession, error) {
	uploadSessions, err := NewUploadSessionsFromRange(size, microsoftGraphUploadSession)
	if err != nil {
		return nil, err
	}
	return &uploadSessions[0], nil
}

func NewUploadSessionsFromRange(size int64, microsoftGraphUploadSession *graphapi.MicrosoftGraphUploadSession) ([]UploadSession, error) {
	uploadSessions := []UploadSession{}
	for _, nextExpectedRange := range microsoftGraphUploadSession.NextExpectedRanges {
		rangeStr := strings.Split(nextExpectedRange, "-")
		from, err := strconv.ParseInt(rangeStr[0], 10, 64)
		if err != nil {
			return nil, err
		}
		to := int64(-1)
		if rangeStr[1] != "" {
			to, err = strconv.ParseInt(rangeStr[1], 10, 64)
			if err != nil {
				return nil, err
			}
		}
		uploadSessionDescription := &UploadSessionDescription{
			Status:        "Wait",
			ContentLength: size,
			ContentRange: UploadSessionDescriptionContentRange{
				Type: "bytes",
				From: from,
				To:   to,
			},
		}
		uploadSessionDescription.ContentRange.To = uploadSessionDescription.SetContentRangTo()
		uploadSession := &UploadSession{
			UploadSessionDescription: uploadSessionDescription,
			UploadSessionReference:   microsoftGraphUploadSession,
		}
		uploadSessions = append(uploadSessions, *uploadSession)
	}
	return uploadSessions, nil
}

func (us *UploadSession) Put(url string, payload io.Reader) (*graphapi.MicrosoftGraphUploadSession, error) {
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, err
	}
	usd := us.UploadSessionDescription
	log.Println("Content-Length: "+usd.GetContentChunkSize(), "Content-Range: "+usd.GetContentRange())
	req.Header.Add("Content-Length", usd.GetContentChunkSize())
	req.Header.Add("Content-Range", usd.GetContentRange())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	if resp.StatusCode < http.StatusBadRequest {
		microsoftGraphUploadSession := graphapi.MicrosoftGraphUploadSession{}
		if err := json.Unmarshal(body, &microsoftGraphUploadSession); err != nil {
			log.Println(err)
		}
		us.UploadSessionDescription.Status = "Finished"
		return &microsoftGraphUploadSession, nil
	}
	if resp.StatusCode < http.StatusInternalServerError {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}
	return nil, errors.New(http.StatusText(resp.StatusCode))
}

func (u *Uploader) Copy(api MicrosoftGraphAPI) {

}

func (u *Uploader) Close(api MicrosoftGraphAPI) {

}

func UseMicrosoftGraphAPIMeDrivecreateUploadSessionPath(str string) string {
	return "/me" + str + ":/createUploadSession"
}

func (usd *UploadSessionDescription) SetContentRangTo() int64 {
	cr := usd.ContentRange
	if cr.To < 0 {
		cr.To = usd.ContentLength - 1
	}
	if cr.To-cr.From > 62914560 {
		cr.To = cr.From + 62914560 - 1
	}
	if cr.To >= usd.ContentLength {
		cr.To = usd.ContentLength - 1
	}
	return cr.To
}

func (usd *UploadSessionDescription) GetContentChunkSizeInt64() int64 {
	cr := usd.ContentRange
	return cr.To - cr.From + 1
}

func (usd *UploadSessionDescription) GetContentChunkSize() string {
	return fmt.Sprintf("%d", usd.GetContentChunkSizeInt64())
}

func (usd *UploadSessionDescription) GetContentLength() string {
	return fmt.Sprintf("%d", usd.ContentLength)
}

func (usd *UploadSessionDescription) GetContentRange() string {
	cr := usd.ContentRange
	return fmt.Sprintf("%s %d-%d/%d", cr.Type, cr.From, cr.To, usd.ContentLength)
}
