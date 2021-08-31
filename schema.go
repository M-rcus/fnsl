package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type API_Subscriptions struct {
	Success  bool                       `json:"success"`
	Response API_Subscriptions_Response `json:"response"`
}

type API_Account struct {
	Success  bool                 `json:"success"`
	Response []API_Common_Account `json:"response"`
}

type API_Subscriptions_Response struct {
	Subscriptions []API_Subscriptions_Response_Subscription `json:"subscriptions"`
}

type API_Subscriptions_Response_Subscription struct {
	AccountId string `json:"accountId"`
}

type API_Timeline struct {
	Success  bool                  `json:"success"`
	Response API_Timeline_Response `json:"response"`
}

type API_Timeline_Response struct {
	Posts               []API_Timeline_Response_Post               `json:"posts"`
	AccountMediaBundles []API_Timeline_Response_AccountMediaBundle `json:"accountMediaBundles"`
	AccountMedia        []API_Timeline_Response_AccountMedia       `json:"accountMedia"`
	Accounts            []API_Common_Account                       `json:"accounts"`
	//AggregatedPosts     []interface{}                               `json:"aggregatedPosts"`
	//Stories             []interface{}                               `json:"stories"`
}

type API_Common_Account struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	DisplayName   string `json:"displayName"`
	ProfileAccess bool   `json:"profileAccess"`
}

type API_Timeline_Response_Post struct {
	Id          string                                  `json:"id"`
	AccountId   string                                  `json:"accountId"`
	Content     string                                  `json:"content"`
	CreatedAt   int64                                   `json:"createdAt"`
	Attachments []API_Timeline_Response_Post_Attachment `json:"attachments"`
}

type API_Timeline_Response_Post_Attachment struct {
	Pos         string `json:"pos"`
	ContentType int    `json:"contentType"`
	ContentId   string `json:"contentId"`
}

type API_Timeline_Response_AccountMedia struct {
	Id        string            `json:"id"`
	AccountId string            `json:"accountId"`
	CreatedAt int64             `json:"createdAt"`
	Media     *API_Common_Media `json:"media"`
	Access    bool              `json:"access"`
}

type API_Common_Media struct {
	Id        string                     `json:"id"`
	Type      int                        `json:"type"`
	Status    int                        `json:"status"`
	AccountId string                     `json:"accountId"`
	MimeType  string                     `json:"mimetype"`
	Filename  string                     `json:"filename"`
	Width     int                        `json:"width"`
	Height    int                        `json:"height"`
	CreatedAt int64                      `json:"createdAt"`
	Variants  []API_Common_Media_Variant `json:"variants"`
	Locations []API_Common_Location      `json:"locations"`
}

func (m *API_Common_Media) IsPhoto() bool {
	return m.Type == API_Common_ContentType_Photo || m.Type == API_Common_ContentType_Thumbnail
}

func (m *API_Common_Media) IsVideo() bool {
	return m.Type == API_Common_ContentType_Video
}

var isNumeric = regexp.MustCompile("^[1-9][0-9]+$")

func (m *API_Common_Media) IsOk() bool {
	if m.Id == "" || m.Type == 0 || m.Status != 1 || m.MimeType == "" || m.CreatedAt < 1 || len(m.Locations) == 0 {
		return false
	}
	if !isNumeric.MatchString(m.Id) || !isNumeric.MatchString(m.AccountId) {
		return false
	}

	return true
}

func (m *API_Common_Media) LocalFilename() string {
	t := time.Unix(m.CreatedAt, 0).UTC()
	fn := t.Format("20060102-150405") + "_" + m.Id
	switch m.MimeType {
	case "image/png":
		fn += ".png"
	case "image/jpg", "image/jpeg":
		fn += ".jpg"
	case "video/mp4":
		fn += ".mp4"
	}
	return fn
}

func (m *API_Common_Media) LocalExists(outDir string) (bool, error) {
	fn := m.LocalFilename()
	_, err := os.Stat(filepath.Join(outDir, fn))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *API_Common_Media) LocalDownload(outDir string) error {
	loc := m.Locations[0]
	if loc.Location == "" {
		return errors.New("location is empty")
	}

	res, err := http.Get(loc.Location)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("invalid response status code: %d", res.StatusCode)
	}

	fn := filepath.Join(outDir, m.LocalFilename())
	outFile, err := os.Create(fn + ".part")
	if err != nil {
		return err
	}
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}

	err = os.Rename(fn+".part", fn)
	if err != nil {
		return err
	}

	return nil
}

type API_Common_Media_Variant struct {
	Id        string                `json:"id"`
	Type      int                   `json:"type"`
	Status    int                   `json:"status"`
	MimeType  string                `json:"mimetype"`
	Filename  string                `json:"filename"`
	Width     int                   `json:"width"`
	Height    int                   `json:"height"`
	Locations []API_Common_Location `json:"locations"`
}

type API_Common_Location struct {
	LocationId string `json:"locationId"`
	Location   string `json:"location"`
}

type API_Timeline_Response_AccountMediaBundle struct {
	Id              string                                                   `json:"id"`
	AccountId       string                                                   `json:"accountId"`
	CreatedAt       int64                                                    `json:"createdAt"`
	AccountMediaIds []string                                                 `json:"accountMediaIds"`
	BundleContent   []API_Timeline_Response_AccountMediaBundle_BundleContent `json:"bundleContent"`
	Access          bool                                                     `json:"access"`
}

type API_Timeline_Response_AccountMediaBundle_BundleContent struct {
	Pos            int    `json:"pos"`
	AccountMediaId string `json:"accountMediaId"`
}

const (
	API_Common_ContentType_Photo     = 1
	API_Common_ContentType_Video     = 2
	API_Common_ContentType_Thumbnail = 3
	API_Common_ContentType_TipGoal   = 7100
)
