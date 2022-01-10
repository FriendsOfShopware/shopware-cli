package account_api

import (
	"encoding/json"
	"fmt"
)

type producerEndpoint struct {
	c          Client
	producerId int
}

func (e producerEndpoint) GetId() int {
	return e.producerId
}

func (c Client) Producer() (*producerEndpoint, error) {
	r, err := c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/companies/%d/allocations", ApiUrl, c.GetActiveCompanyId()), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var allocation companyAllocation
	if err := json.Unmarshal(body, &allocation); err != nil {
		return nil, fmt.Errorf("my_profile: %v", err)
	}

	if !allocation.IsProducer {
		return nil, fmt.Errorf("this company is not unlocked as producer")
	}

	return &producerEndpoint{producerId: allocation.ProducerId, c: c}, nil
}

type companyAllocation struct {
	HasShops          bool `json:"hasShops"`
	HasCommercialShop bool `json:"hasCommercialShop"`
	IsEducationMember bool `json:"isEducationMember"`
	IsPartner         bool `json:"isPartner"`
	IsProducer        bool `json:"isProducer"`
	ProducerId        int  `json:"producerId"`
}

func (e producerEndpoint) Profile() (*producer, error) {
	// Fetch the producer
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/producers?companyId=%d", ApiUrl, e.c.GetActiveCompanyId()), nil)
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var producers []producer
	if err := json.Unmarshal(body, &producers); err != nil {
		return nil, fmt.Errorf("my_profile: %v", err)
	}

	for _, profile := range producers {
		return &profile, nil
	}

	return nil, fmt.Errorf("cannot find a profile")
}

type producer struct {
	Id       int    `json:"id"`
	Prefix   string `json:"prefix"`
	Contract struct {
		Id   int    `json:"id"`
		Path string `json:"path"`
	} `json:"contract"`
	Name    string `json:"name"`
	Details []struct {
		Id     int `json:"id"`
		Locale struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"locale"`
		Description string `json:"description"`
	} `json:"details"`
	Website              string `json:"website"`
	Fixed                bool   `json:"fixed"`
	HasCancelledContract bool   `json:"hasCancelledContract"`
	IconPath             string `json:"iconPath"`
	IconIsSet            bool   `json:"iconIsSet"`
	ShopwareId           string `json:"shopwareId"`
	UserId               int    `json:"userId"`
	CompanyId            int    `json:"companyId"`
	CompanyName          string `json:"companyName"`
	SaleMail             string `json:"saleMail"`
	SupportMail          string `json:"supportMail"`
	RatingMail           string `json:"ratingMail"`
	SupportedLanguages   []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"supportedLanguages"`
	IconUrl           string      `json:"iconUrl"`
	CancelledContract interface{} `json:"cancelledContract"`
}
