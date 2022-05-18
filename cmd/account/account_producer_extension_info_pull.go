package account

import (
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var accountCompanyProducerExtensionInfoPullCmd = &cobra.Command{
	Use:   "pull [path]",
	Short: "Generates local store configuration from account data",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])

		if err != nil {
			return errors.Wrap(err, "cannot open file")
		}

		zipExt, err := extension.GetExtensionByFolder(path)

		if err != nil {
			return errors.Wrap(err, "cannot open extension")
		}

		zipName, err := zipExt.GetName()

		if err != nil {
			return errors.Wrap(err, "cannot get extension name")
		}

		p, err := services.AccountClient.Producer()

		if err != nil {
			return errors.Wrap(err, "cannot get producer endpoint")
		}

		storeExt, err := p.GetExtensionByName(zipName)

		if err != nil {
			return errors.Wrap(err, "cannot get store extension")
		}

		resourcesFolder := fmt.Sprintf("%s/src/Resources/store/", zipExt.GetPath())
		categoryList := make([]string, 0)
		availabilities := make([]string, 0)
		localizations := make([]string, 0)
		tagsDE := make([]string, 0)
		tagsEN := make([]string, 0)
		videosDE := make([]string, 0)
		videosEN := make([]string, 0)
		highlightsDE := make([]string, 0)
		highlightsEN := make([]string, 0)
		featuresDE := make([]string, 0)
		featuresEN := make([]string, 0)
		faqDE := make([]extension.ConfigStoreFaq, 0)
		faqEN := make([]extension.ConfigStoreFaq, 0)
		images := make([]extension.ConfigStoreImage, 0)

		if _, err := os.Stat(resourcesFolder); os.IsNotExist(err) {
			err = os.MkdirAll(resourcesFolder, os.ModePerm)

			if err != nil {
				return errors.Wrap(err, "cannot create file")
			}
		}

		var iconConfigPath *string

		if len(storeExt.IconURL) > 0 {
			icon := "src/Resources/store/icon.png"
			iconConfigPath = &icon
			err := downloadFileTo(storeExt.IconURL, fmt.Sprintf("%s/icon.png", resourcesFolder))
			if err != nil {
				return errors.Wrap(err, "cannot download file")
			}
		}

		for _, category := range storeExt.Categories {
			categoryList = append(categoryList, category.Name)
		}

		for _, localization := range storeExt.Localizations {
			localizations = append(localizations, localization.Name)
		}

		for _, a := range storeExt.StoreAvailabilities {
			availabilities = append(availabilities, a.Name)
		}

		storeImages, err := p.GetExtensionImages(storeExt.Id)

		if err != nil {
			return errors.Wrap(err, "cannot get extension images")
		}

		for i, image := range storeImages {
			imagePath := fmt.Sprintf("src/Resources/store/img-%d.png", i)
			err := downloadFileTo(image.RemoteLink, fmt.Sprintf("%s/%s", zipExt.GetPath(), imagePath))
			if err != nil {
				return errors.Wrap(err, "cannot download file")
			}

			images = append(images, extension.ConfigStoreImage{
				File:     imagePath,
				Preview:  extension.ConfigStoreImagePreview{German: image.Details[0].Preview, English: image.Details[1].Preview},
				Activate: extension.ConfigStoreImageActivate{German: image.Details[0].Activated, English: image.Details[1].Activated},
				Priority: image.Priority,
			})
		}

		for _, info := range storeExt.Infos {
			language := info.Locale.Name[0:2]

			if language == "de" {
				for _, element := range info.Tags {
					tagsDE = append(tagsDE, element.Name)
				}

				for _, element := range info.Videos {
					videosDE = append(videosDE, element.URL)
				}

				highlightsDE = append(highlightsDE, strings.Split(info.Highlights, "\n")...)
				featuresDE = append(featuresDE, strings.Split(info.Features, "\n")...)

				for _, element := range info.Faqs {
					faqDE = append(faqDE, extension.ConfigStoreFaq{Question: element.Question, Answer: element.Answer})
				}
			} else {
				for _, element := range info.Tags {
					tagsEN = append(tagsEN, element.Name)
				}

				for _, element := range info.Videos {
					videosEN = append(videosEN, element.URL)
				}

				highlightsEN = append(highlightsEN, strings.Split(info.Highlights, "\n")...)
				featuresEN = append(featuresEN, strings.Split(info.Features, "\n")...)

				for _, element := range info.Faqs {
					faqEN = append(faqEN, extension.ConfigStoreFaq{Question: element.Question, Answer: element.Answer})
				}
			}
		}

		newCfg := extension.Config{Store: extension.ConfigStore{
			Icon:                                iconConfigPath,
			DefaultLocale:                       &storeExt.StandardLocale.Name,
			Type:                                &storeExt.ProductType.Name,
			AutomaticBugfixVersionCompatibility: &storeExt.AutomaticBugfixVersionCompatibility,
			Availabilities:                      &availabilities,
			Localizations:                       &localizations,
			Description:                         extension.ConfigTranslatedString{German: &storeExt.Infos[0].Description, English: &storeExt.Infos[1].Description},
			InstallationManual:                  extension.ConfigTranslatedString{German: &storeExt.Infos[0].InstallationManual, English: &storeExt.Infos[1].InstallationManual},
			Categories:                          &categoryList,
			Tags:                                extension.ConfigTranslatedStringList{German: &tagsDE, English: &tagsEN},
			Videos:                              extension.ConfigTranslatedStringList{German: &videosDE, English: &videosEN},
			Highlights:                          extension.ConfigTranslatedStringList{German: &highlightsDE, English: &highlightsEN},
			Features:                            extension.ConfigTranslatedStringList{German: &featuresDE, English: &featuresEN},
			Faq:                                 extension.ConfigStoreTranslatedFaq{German: &faqDE, English: &faqEN},
			Images:                              &images,
		}}

		content, err := yaml.Marshal(newCfg)

		if err != nil {
			return errors.Wrap(err, "cannot encode yaml")
		}

		extCfgFile := fmt.Sprintf("%s/%s", zipExt.GetPath(), ".shopware-extension.yml")
		err = ioutil.WriteFile(extCfgFile, content, os.ModePerm)

		if err != nil {
			return errors.Wrap(err, "cannot save file")
		}

		log.Infof("Files has been written to the given extension folder")

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionInfoCmd.AddCommand(accountCompanyProducerExtensionInfoPullCmd)
}

func downloadFileTo(url string, target string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil) //nolint:noctx
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "download file")
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read file body")
	}

	err = ioutil.WriteFile(target, content, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "write to file")
	}

	return nil
}
