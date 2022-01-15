package account_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"net/url"
	"strconv"
	"strings"
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
		return nil, fmt.Errorf("producer.profile: %v", err)
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

type ListExtensionCriteria struct {
	Limit         int    `schema:"limit,omitempty"`
	Offset        int    `schema:"offset,omitempty"`
	OrderBy       string `schema:"orderBy,omitempty"`
	OrderSequence string `schema:"orderSequence,omitempty"`
	Search        string `schema:"search,omitempty"`
}

func (e producerEndpoint) Extensions(criteria *ListExtensionCriteria) ([]extension, error) {
	encoder := schema.NewEncoder()
	form := url.Values{}
	form.Set("producerId", strconv.FormatInt(int64(e.GetId()), 10))
	err := encoder.Encode(criteria, form)

	if err != nil {
		return nil, fmt.Errorf("list_extensions: %v", err)
	}

	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins?%s", ApiUrl, form.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var extensions []extension
	if err := json.Unmarshal(body, &extensions); err != nil {
		return nil, fmt.Errorf("list_extensions: %v", err)
	}

	return extensions, nil
}

func (e producerEndpoint) GetExtensionByName(name string) (*extension, error) {
	criteria := ListExtensionCriteria{
		Search: name,
	}

	extensions, err := e.Extensions(&criteria)

	if err != nil {
		return nil, err
	}

	for _, ext := range extensions {
		if strings.ToLower(ext.Name) == strings.ToLower(name) {
			return e.GetExtensionById(ext.Id)
		}
	}

	return nil, fmt.Errorf("cannot find extension by name %s", name)
}

func (e producerEndpoint) GetExtensionById(id int) (*extension, error) {
	// Create it
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/plugins/%d", ApiUrl, id), nil)

	if err != nil {
		return nil, fmt.Errorf("GetExtensionById: %v", err)
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, fmt.Errorf("GetExtensionById: %v", err)
	}

	var extension extension
	if err := json.Unmarshal(body, &extension); err != nil {
		return nil, fmt.Errorf("GetExtensionById: %v", err)
	}

	return &extension, nil
}

type extension struct {
	Id       int `json:"id"`
	Producer struct {
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
	} `json:"producer"`
	Type struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"type"`
	Name            string `json:"name"`
	Code            string `json:"code"`
	ModuleKey       string `json:"moduleKey"`
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
	StandardLocale struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"standardLocale"`
	License struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"license"`
	Infos []struct {
		Id     int `json:"id"`
		Locale struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"locale"`
		Name               string        `json:"name"`
		Description        string        `json:"description"`
		InstallationManual string        `json:"installationManual"`
		ShortDescription   string        `json:"shortDescription"`
		Highlights         string        `json:"highlights"`
		Features           string        `json:"features"`
		Tags               []interface{} `json:"tags"`
		Videos             []interface{} `json:"videos"`
		Faqs               []interface{} `json:"faqs"`
	} `json:"infos"`
	PriceModels         []interface{} `json:"priceModels"`
	Variants            []interface{} `json:"variants"`
	StoreAvailabilities []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"storeAvailabilities"`
	Categories                          []interface{} `json:"categories"`
	Addons                              []interface{} `json:"addons"`
	LastChange                          string        `json:"lastChange"`
	CreationDate                        string        `json:"creationDate"`
	Support                             bool          `json:"support"`
	SupportOnlyCommercial               bool          `json:"supportOnlyCommercial"`
	IconPath                            string        `json:"iconPath"`
	IconIsSet                           bool          `json:"iconIsSet"`
	ExamplePageUrl                      string        `json:"examplePageUrl"`
	Demos                               []interface{} `json:"demos"`
	Localizations                       []interface{} `json:"localizations"`
	LatestBinary                        interface{}   `json:"latestBinary"`
	MigrationSupport                    bool          `json:"migrationSupport"`
	AutomaticBugfixVersionCompatibility bool          `json:"automaticBugfixVersionCompatibility"`
	HiddenInStore                       bool          `json:"hiddenInStore"`
	Certification                       interface{}   `json:"certification"`
	ProductType                         interface{}   `json:"productType"`
	Status                              struct {
		Name string `json:"name"`
	} `json:"status"`
	MinimumMarketingSoftwareVersion       interface{} `json:"minimumMarketingSoftwareVersion"`
	IsSubscriptionEnabled                 bool        `json:"isSubscriptionEnabled"`
	ReleaseDate                           interface{} `json:"releaseDate"`
	PlannedReleaseDate                    interface{} `json:"plannedReleaseDate"`
	LastBusinessModelChangeDate           interface{} `json:"lastBusinessModelChangeDate"`
	IsSW5Compatible                       bool        `json:"isSW5Compatible"`
	Subprocessors                         interface{} `json:"subprocessors"`
	PluginTestingInstanceDisabled         bool        `json:"pluginTestingInstanceDisabled"`
	IconUrl                               string      `json:"iconUrl"`
	Pictures                              string      `json:"pictures"`
	HasPictures                           bool        `json:"hasPictures"`
	Comments                              string      `json:"comments"`
	Reviews                               string      `json:"reviews"`
	IsPremiumPlugin                       bool        `json:"isPremiumPlugin"`
	IsAdvancedFeature                     bool        `json:"isAdvancedFeature"`
	IsEnterpriseAccelerator               bool        `json:"isEnterpriseAccelerator"`
	IsSW6EnterpriseFeature                bool        `json:"isSW6EnterpriseFeature"`
	IsSW6ProfessionalEditionFeature       bool        `json:"isSW6ProfessionalEditionFeature"`
	Binaries                              interface{} `json:"binaries"`
	Predecessor                           interface{} `json:"predecessor"`
	Successor                             interface{} `json:"successor"`
	IsCompatibleWithLatestShopwareVersion bool        `json:"isCompatibleWithLatestShopwareVersion"`
	PluginPreview                         interface{} `json:"pluginPreview"`
}

type CreateExtensionRequest struct {
	Name       string `json:"name,omitempty"`
	Generation struct {
		Name string `json:"name"`
	} `json:"generation"`
	ProducerId int `json:"producerId"`
}

const (
	GenerationClassic  = "classic"
	GenerationThemes   = "themes"
	GenerationApps     = "apps"
	GenerationPlatform = "platform"
)

func (e producerEndpoint) CreateExtension(newExtension CreateExtensionRequest) (*extension, error) {
	requestBody, err := json.Marshal(newExtension)

	if err != nil {
		return nil, err
	}

	// Create it
	r, err := e.c.NewAuthenticatedRequest("POST", fmt.Sprintf("%s/plugins", ApiUrl), bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, err
	}

	var extension extension
	if err := json.Unmarshal(body, &extension); err != nil {
		return nil, fmt.Errorf("create_extension: %v", err)
	}

	extension.Name = newExtension.Name

	// Patch the name
	err = e.UpdateExtension(&extension)

	if err != nil {
		return nil, err
	}

	return &extension, nil
}

func (e producerEndpoint) UpdateExtension(extension *extension) error {
	requestBody, err := json.Marshal(extension)

	if err != nil {
		return err
	}

	// Patch the name
	r, err := e.c.NewAuthenticatedRequest("PUT", fmt.Sprintf("%s/plugins/%d", ApiUrl, extension.Id), bytes.NewBuffer(requestBody))

	if err != nil {
		return err
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) DeleteExtension(id int) error {
	r, err := e.c.NewAuthenticatedRequest("DELETE", fmt.Sprintf("%s/plugins/%d", ApiUrl, id), nil)

	if err != nil {
		return err
	}

	_, err = e.c.doRequest(r)

	return err
}

func (e producerEndpoint) GetSoftwareVersions(generation string) (*[]softwareVersion, error) {
	r, err := e.c.NewAuthenticatedRequest("GET", fmt.Sprintf("%s/pluginstatics/softwareVersions?filter=[{\"property\":\"pluginGeneration\",\"value\":\"%s\"}]", ApiUrl, generation), nil)

	if err != nil {
		return nil, fmt.Errorf("shopware_versions: %v", err)
	}

	body, err := e.c.doRequest(r)

	if err != nil {
		return nil, fmt.Errorf("shopware_versions: %v", err)
	}

	var versions []softwareVersion

	err = json.Unmarshal(body, &versions)

	if err != nil {
		return nil, fmt.Errorf("shopware_versions: %v", err)
	}

	return &versions, nil
}

type softwareVersion struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Parent      interface{} `json:"parent"`
	Selectable  bool        `json:"selectable"`
	Major       *string     `json:"major"`
	ReleaseDate *string     `json:"releaseDate"`
	Public      bool        `json:"public"`
}
