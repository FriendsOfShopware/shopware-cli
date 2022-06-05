package account_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/version"
	"github.com/microcosm-cc/bluemonday"
)

type SoftwareVersionList []SoftwareVersion

type ExtensionBinary struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	RemoteLink string `json:"remoteLink"`
	Version    string `json:"version"`
	Status     struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"status"`
	CompatibleSoftwareVersions SoftwareVersionList `json:"compatibleSoftwareVersions"`
	Changelogs                 []struct {
		Id     int `json:"id"`
		Locale struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"locale"`
		Text string `json:"text"`
	} `json:"changelogs"`
	CreationDate   string `json:"creationDate"`
	LastChangeDate string `json:"lastChangeDate"`
	Archives       []struct {
		Id                   int         `json:"id"`
		RemoteLink           string      `json:"remoteLink"`
		ShopwareMajorVersion interface{} `json:"shopwareMajorVersion"`
		IoncubeEncrypted     bool        `json:"ioncubeEncrypted"`
		ManifestRemoteLink   interface{} `json:"manifestRemoteLink"`
	} `json:"archives"`
	IonCubeEncrypted            bool `json:"ionCubeEncrypted"`
	LicenseCheckRequired        bool `json:"licenseCheckRequired"`
	HasActiveCodeReviewWarnings bool `json:"hasActiveCodeReviewWarnings"`
}

func (e producerEndpoint) GetExtensionBinaries(extensionId int) ([]*ExtensionBinary, error) {
	errorFormat := "GetExtensionBinaries: %v"

	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d/binaries", ApiUrl, extensionId), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var binaries []*ExtensionBinary
	if err := json.Unmarshal(body, &binaries); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return binaries, nil
}

func (e producerEndpoint) UpdateExtensionBinaryInfo(extensionId int, binary ExtensionBinary) error {
	errorFormat := "UpdateExtensionBinaryInfo: %v"

	content, err := json.Marshal(binary)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("PUT", fmt.Sprintf("%s/plugins/%d/binaries/%d", ApiUrl, extensionId, binary.Id), bytes.NewReader(content)) //nolint:noctx
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) CreateExtensionBinaryFile(extensionId int, zipPath string) (*ExtensionBinary, error) {
	errorFormat := "CreateExtensionBinaryFile: %v"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(zipPath))
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	zipFile, err := os.Open(zipPath)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	_, err = io.Copy(fileWritter, zipFile)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/binaries", ApiUrl, extensionId), &b)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	content, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	// For some reasons this API responses a array of binaries
	var binary []*ExtensionBinary
	if err := json.Unmarshal(content, &binary); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return binary[0], nil
}

func (e producerEndpoint) UpdateExtensionBinaryFile(extensionId, binaryId int, zipPath string) error {
	errorFormat := "UpdateExtensionBinaryFile: %v"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(zipPath))
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	zipFile, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	if _, err = io.Copy(fileWritter, zipFile); err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/binaries/%d/file", ApiUrl, extensionId, binaryId), &b) //nolint:noctx
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) UpdateExtensionIcon(extensionId int, iconFile string) error {
	errorFormat := "UpdateExtensionIcon: %v"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(iconFile))
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	zipFile, err := os.Open(iconFile)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	if _, err = io.Copy(fileWritter, zipFile); err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/icon", ApiUrl, extensionId), &b)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	_, err = e.c.doRequest(r)

	return err
}

type ExtensionImage struct {
	Id         int    `json:"id"`
	RemoteLink string `json:"remoteLink"`
	Details    []struct {
		Id        int    `json:"id"`
		Preview   bool   `json:"preview"`
		Activated bool   `json:"activated"`
		Caption   string `json:"caption"`
		Locale    Locale `json:"locale"`
	} `json:"details"`
	Priority int `json:"priority"`
}

func (e producerEndpoint) GetExtensionImages(extensionId int) ([]*ExtensionImage, error) {
	errorFormat := "GetExtensionImages: %v"

	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d/pictures", ApiUrl, extensionId), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var images []*ExtensionImage
	if err := json.Unmarshal(body, &images); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return images, nil
}

func (e producerEndpoint) DeleteExtensionImages(extensionId, imageId int) error {
	errorFormat := "DeleteExtensionImages: %v"

	r, err := e.c.NewAuthenticatedRequest("DELETE", fmt.Sprintf("%s/plugins/%d/pictures/%d", ApiUrl, extensionId, imageId), nil)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) UpdateExtensionImage(extensionId int, image *ExtensionImage) error {
	errorFormat := "UpdateExtensionImage: %v"

	content, err := json.Marshal(image)
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("PUT", fmt.Sprintf("%s/plugins/%d/pictures/%d", ApiUrl, extensionId, image.Id), bytes.NewReader(content)) //nolint:noctx
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) AddExtensionImage(extensionId int, file string) (*ExtensionImage, error) {
	errorFormat := "AddExtensionImage: %v"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(file))
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	zipFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	if _, err = io.Copy(fileWritter, zipFile); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/pictures", ApiUrl, extensionId), &b) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var list []*ExtensionImage

	if err = json.Unmarshal(body, &list); err != nil {
		return nil, fmt.Errorf("AddExtensionImage: %v", err)
	}

	return list[0], nil
}

func (e producerEndpoint) TriggerCodeReview(extensionId int) error {
	errorFormat := "TriggerCodeReview: %v"

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/reviews", ApiUrl, extensionId), nil) //nolint:noctx
	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) GetBinaryReviewResults(extensionId, binaryId int) ([]BinaryReviewResult, error) {
	errorFormat := "GetBinaryReviewResults: %v"

	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d/binaries/%d/checkresults", ApiUrl, extensionId, binaryId), nil) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	body, err := e.c.doRequest(r)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	var results []BinaryReviewResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return results, nil
}

type BinaryReviewResult struct {
	Id       int `json:"id"`
	BinaryId int `json:"binaryId"`
	Type     struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"type"`
	Message         string `json:"message"`
	CreationDate    string `json:"creationDate"`
	SubCheckResults []struct {
		SubCheck    string `json:"subCheck"`
		Status      string `json:"status"`
		Passed      bool   `json:"passed"`
		Message     string `json:"message"`
		HasWarnings bool   `json:"hasWarnings"`
	} `json:"subCheckResults"`
}

func (review BinaryReviewResult) HasPassed() bool {
	return review.Type.Id == 3 || review.Type.Name == "automaticcodereviewsucceeded"
}

func (review BinaryReviewResult) HasWarnings() bool {
	for _, result := range review.SubCheckResults {
		if result.HasWarnings {
			return true
		}
	}

	return false
}

func (review BinaryReviewResult) IsPending() bool {
	return review.Type.Id == 4
}

func (review BinaryReviewResult) GetSummary() string {
	message := ""

	p := bluemonday.NewPolicy()

	for _, result := range review.SubCheckResults {
		if result.Passed && !result.HasWarnings {
			continue
		}

		message += fmt.Sprintf("=== %s ===\n", result.SubCheck)
		message += fmt.Sprintf("%s\n\n", p.Sanitize(result.Message))
	}

	return message
}

func (list SoftwareVersionList) FilterOnVersion(constriant *version.Constraints) SoftwareVersionList {
	newList := make(SoftwareVersionList, 0)

	for _, swVersion := range list {
		if !swVersion.Selectable {
			continue
		}

		v, err := version.NewVersion(swVersion.Name)
		if err != nil {
			continue
		}

		if constriant.Check(v) {
			newList = append(newList, swVersion)
		}
	}

	return newList
}
