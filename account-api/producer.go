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

func (e producerEndpoint) Extensions() ([]listingExtension, error) {
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins?producerId=%d&limit=100&orderBy=name&orderSequence=asc", ApiUrl, e.GetId()), nil)
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var extensions []listingExtension
	if err := json.Unmarshal(body, &extensions); err != nil {
		return nil, fmt.Errorf("list_extensions: %v", err)
	}

	return extensions, nil
}

type listingExtension struct {
	Id   int `json:"id"`
	Type struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"type"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	LifecycleStatus struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"lifecycleStatus"`
	Generation struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"generation"`
	ActivationStatus struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"activationStatus"`
	ApprovalStatus struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"approvalStatus"`
	Infos []struct {
		Id     int `json:"id"`
		Locale struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"locale"`
		Name               string `json:"name"`
		Description        string `json:"description"`
		InstallationManual string `json:"installationManual"`
		ShortDescription   string `json:"shortDescription"`
		Highlights         string `json:"highlights"`
		Features           string `json:"features"`
		Tags               []struct {
			Id     int `json:"id"`
			Locale struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"locale"`
			Name     string `json:"name"`
			Internal bool   `json:"internal"`
		} `json:"tags"`
		Videos []interface{} `json:"videos"`
		Faqs   []interface{} `json:"faqs"`
	} `json:"infos"`
	Variants []struct {
		Id         int `json:"id"`
		PriceModel struct {
			Id          int    `json:"id"`
			Type        string `json:"type"`
			BookingKey  string `json:"bookingKey"`
			BookingText string `json:"bookingText"`
		} `json:"priceModel"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"variants"`
	StoreAvailabilities []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"storeAvailabilities"`
	Addons        []interface{} `json:"addons"`
	CreationDate  string        `json:"creationDate"`
	Certification struct {
		PluginId     int    `json:"pluginId"`
		CreationDate string `json:"creationDate"`
		BronzeResult struct {
			Items []struct {
				Name      string      `json:"name"`
				Label     string      `json:"label"`
				Actual    interface{} `json:"actual"`
				Expected  interface{} `json:"expected"`
				Fulfilled bool        `json:"fulfilled"`
			} `json:"items"`
			Fulfilled bool `json:"fulfilled"`
		} `json:"bronzeResult"`
		SilverResult struct {
			Items []struct {
				Name      string      `json:"name"`
				Label     string      `json:"label"`
				Actual    interface{} `json:"actual"`
				Expected  interface{} `json:"expected"`
				Fulfilled bool        `json:"fulfilled"`
			} `json:"items"`
			Fulfilled bool `json:"fulfilled"`
		} `json:"silverResult"`
		GoldResult struct {
			Items []struct {
				Name      string      `json:"name"`
				Label     string      `json:"label"`
				Actual    interface{} `json:"actual"`
				Expected  interface{} `json:"expected"`
				Fulfilled bool        `json:"fulfilled"`
			} `json:"items"`
			Fulfilled bool `json:"fulfilled"`
		} `json:"goldResult"`
		Type struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"type"`
	} `json:"certification"`
	ProductType struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		MainCategory struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Parent      bool   `json:"parent"`
			Position    int    `json:"position"`
			Public      bool   `json:"public"`
			Visible     bool   `json:"visible"`
			Suggested   bool   `json:"suggested"`
			Applicable  bool   `json:"applicable"`
			Details     []struct {
				Id          int    `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Locale      struct {
					Id   int    `json:"id"`
					Name string `json:"name"`
				} `json:"locale"`
			} `json:"details"`
			Active bool `json:"active"`
		} `json:"mainCategory"`
	} `json:"productType"`
	Status struct {
		Name string `json:"name"`
	} `json:"status"`
	PlannedReleaseDate                    interface{} `json:"plannedReleaseDate"`
	Successor                             interface{} `json:"successor"`
	IsSW5Compatible                       bool        `json:"isSW5Compatible"`
	IsCompatibleWithLatestShopwareVersion bool        `json:"isCompatibleWithLatestShopwareVersion"`
	AutomaticBugfixVersionCompatibility   bool        `json:"automaticBugfixVersionCompatibility"`
	ReleaseDate                           *struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"releaseDate"`
	PluginTestingInstanceCreated bool `json:"pluginTestingInstanceCreated"`
}
