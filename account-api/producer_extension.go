package account_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d/binaries", ApiUrl, extensionId), nil)

	if err != nil {
		return nil, fmt.Errorf("GetExtensionBinaries: %v", err)
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, fmt.Errorf("GetExtensionBinaries: %v", err)
	}

	var binaries []*ExtensionBinary
	if err := json.Unmarshal(body, &binaries); err != nil {
		return nil, fmt.Errorf("GetExtensionBinaries: %v", err)
	}

	return binaries, nil
}

func (e producerEndpoint) UpdateExtensionBinaryInfo(extensionId int, binary ExtensionBinary) error {
	content, err := json.Marshal(binary)

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryInfo: %v", err)
	}

	r, err := e.c.NewAuthenticatedRequest("PUT", fmt.Sprintf("%s/plugins/%d/binaries/%d", ApiUrl, extensionId, binary.Id), bytes.NewReader(content))

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryInfo: %v", err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) CreateExtensionBinaryFile(extensionId int, zipPath string) (*ExtensionBinary, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(zipPath))

	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	zipFile, err := os.Open(zipPath)

	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	_, err = io.Copy(fileWritter, zipFile)

	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/binaries", ApiUrl, extensionId), &b)

	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	content, err := e.c.doRequest(r)

	if err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	// For some reasons this API responses a array of binaries
	var binary []*ExtensionBinary
	if err := json.Unmarshal(content, &binary); err != nil {
		return nil, fmt.Errorf("CreateExtensionBinaryFile: %v", err)
	}

	return binary[0], nil
}

func (e producerEndpoint) UpdateExtensionBinaryFile(extensionId, binaryId int, zipPath string) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileWritter, err := w.CreateFormFile("file", filepath.Base(zipPath))

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryFile: %v", err)
	}

	zipFile, err := os.Open(zipPath)

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryFile: %v", err)
	}

	_, err = io.Copy(fileWritter, zipFile)

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryFile: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryFile: %v", err)
	}

	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/binaries/%d/file", ApiUrl, extensionId, binaryId), &b)

	if err != nil {
		return fmt.Errorf("UpdateExtensionBinaryFile: %v", err)
	}

	r.Header.Set("content-type", w.FormDataContentType())

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) TriggerCodeReview(extensionId int) error {
	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins/%d/reviews", ApiUrl, extensionId), nil)

	if err != nil {
		return fmt.Errorf("TriggerCodeReview: %v", err)
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) GetBinaryReviewResults(extensionId, binaryId int) ([]BinaryReviewResult, error) {
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d/binaries/%d/checkresults", ApiUrl, extensionId, binaryId), nil)

	if err != nil {
		return nil, fmt.Errorf("GetBinaryReviewResults: %v", err)
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, fmt.Errorf("GetBinaryReviewResults: %v", err)
	}

	var results []BinaryReviewResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("GetBinaryReviewResults: %v", err)
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

		message = message + fmt.Sprintf("=== %s ===\n", result.SubCheck)
		message = message + fmt.Sprintf("%s\n\n", p.Sanitize(result.Message))
	}

	return message
}

func (list SoftwareVersionList) FilterOnVersion(constriant *version.Constraints) SoftwareVersionList {
	newList := make(SoftwareVersionList, 0)

	for _, swVersion := range list {
		if swVersion.Selectable == false {
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
