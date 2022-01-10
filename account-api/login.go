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

	memberships, err := fetchMemberships(token)

	if err != nil {
		return nil, err
	}

	var activeMemberShip membership

	for _, membership := range *memberships {
		if membership.Company.Id == token.UserID {
			activeMemberShip = membership
		}
	}

	client := Client{
		token:            token,
		memberships:      memberships,
		activeMembership: &activeMemberShip,
	}

	return &client, nil
}

func fetchMemberships(token token) (*[]membership, error) {
	r, err := http.NewRequest("GET", fmt.Sprintf("%s/account/%d/memberships", ApiUrl, token.UserAccountID), nil)
	r.Header.Set("x-shopware-token", token.Token)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetchMemberships: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(string(data))
	}

	var companies []membership
	if err := json.Unmarshal(data, &companies); err != nil {
		return nil, fmt.Errorf("fetchMemberships: %v", err)
	}

	return &companies, nil
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

type membership struct {
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
}

type changeMembershipRequest struct {
	SelectedMembership struct {
		Id int `json:"id"`
	} `json:"membership"`
}

func (c Client) ChangeActiveMembership(selected *membership) error {
	s, err := json.Marshal(changeMembershipRequest{SelectedMembership: struct {
		Id int `json:"id"`
	}(struct{ Id int }{Id: selected.Id})})

	if err != nil {
		return fmt.Errorf("ChangeActiveMembership: %v", err)
	}

	r, err := c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/account/%d/memberships/change", ApiUrl, c.GetUserId()), bytes.NewBuffer(s))

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		c.activeMembership = selected

		return nil
	}

	return fmt.Errorf("could not change active membership due http error %d", resp.StatusCode)
}
