package account_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c Client) GetMyProfile() (*myProfile, error) {
	request, err := c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/account/%d", ApiUrl, c.token.UserAccountID), nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(data))
	}

	var profile myProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("my_profile: %v", err)
	}

	return &profile, nil
}

type myProfile struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	CreationDate string `json:"creationDate"`
	Banned       bool   `json:"banned"`
	Verified     bool   `json:"verified"`
	PersonalData struct {
		Id         int `json:"id"`
		Salutation struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"salutation"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Locale    struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"locale"`
	} `json:"personalData"`
	PartnerMarketingOptIn bool `json:"partnerMarketingOptIn"`
	SelectedMembership    struct {
		Id           int    `json:"id"`
		CreationDate string `json:"creationDate"`
		Active       bool   `json:"active"`
		Member       struct {
			Id           int         `json:"id"`
			Email        string      `json:"email"`
			AvatarUrl    interface{} `json:"avatarUrl"`
			PersonalData struct {
				Id         int `json:"id"`
				Salutation struct {
					Id          int    `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"salutation"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    struct {
					Id          int    `json:"id"`
					Name        string `json:"name"`
					Description string `json:"description"`
				} `json:"locale"`
			} `json:"personalData"`
		} `json:"member"`
		Company struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			CustomerNumber int    `json:"customerNumber"`
		} `json:"company"`
		Roles []struct {
			Id           int         `json:"id"`
			Name         string      `json:"name"`
			CreationDate string      `json:"creationDate"`
			Company      interface{} `json:"company"`
			Permissions  []struct {
				Id      int    `json:"id"`
				Context string `json:"context"`
				Name    string `json:"name"`
			} `json:"permissions"`
		} `json:"roles"`
	} `json:"selectedMembership"`
}
