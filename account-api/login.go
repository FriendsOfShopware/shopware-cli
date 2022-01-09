package account_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ApiUrl = "https://api.shopware.com"

func NewApi(request LoginRequest) (*Client, error) {
	s, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	resp, err := http.Post(ApiUrl+"/accesstokens", "application/json", bytes.NewBuffer(s))
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(data))
	}

	var token token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	client := Client{
		token: token,
	}

	return &client, nil
}

type token struct {
	Token         string      `json:"token"`
	Expire        tokenExpire `json:"expire"`
	UserAccountID int         `json:"userAccountId"`
	UserID        int         `json:"userId"`
	LegacyLogin   bool        `json:"legacyLogin"`
}

type tokenExpire struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"shopwareId"`
	Password string `json:"password"`
}
