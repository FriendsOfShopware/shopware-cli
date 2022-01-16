package cmd

import (
	"fmt"
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	accountApi "shopware-cli/account-api"
	"shopware-cli/extension"
	"strings"
)

var accountCompanyProducerExtensionUpdateCmd = &cobra.Command{
	Use:   "update [zip or path]",
	Short: "Update store information of extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountApiByConfig()

		path, err := filepath.Abs(args[0])

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		stat, err := os.Stat(path)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		var zipExt extension.Extension

		if stat.IsDir() {
			zipExt, err = extension.GetExtensionByFolder(path)
		} else {
			zipExt, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		zipName, err := zipExt.GetName()

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		p, err := client.Producer()

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		storeExt, err := p.GetExtensionByName(zipName)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		metadata := zipExt.GetMetaData()

		for _, info := range storeExt.Infos {
			language := info.Locale.Name[0:2]

			if language == "de" {
				info.Name = metadata.Label.German
				info.ShortDescription = metadata.Description.German
			} else {
				info.Name = metadata.Label.English
				info.ShortDescription = metadata.Description.English
			}
		}

		extCfg, err := extension.ReadExtensionConfig(zipExt.GetPath())
		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		info, err := p.GetExtensionGeneralInfo()

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		if extCfg != nil {
			updateStoreInfo(storeExt, extCfg, info)
		}

		err = p.UpdateExtension(storeExt)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		termColor.Green("Store information has been updated")
	},
}

func updateStoreInfo(ext *accountApi.Extension, cfg *extension.Config, info *accountApi.ExtensionGeneralInformation) {
	if cfg.Store.DefaultLocale != nil {
		for _, locale := range info.Locales {
			if locale.Name == *cfg.Store.DefaultLocale {
				ext.StandardLocale = locale
			}
		}
	}

	if cfg.Store.Localizations != nil {
		newLocales := make([]accountApi.Locale, 0)

		for _, locale := range info.Locales {
			for _, configLocale := range *cfg.Store.Localizations {
				if locale.Name == configLocale {
					newLocales = append(newLocales, locale)
				}
			}
		}

		ext.Localizations = newLocales
	}

	if cfg.Store.Availabilities != nil {
		newAvailabilities := make([]accountApi.StoreAvailablity, 0)

		for _, availablity := range info.StoreAvailabilities {
			for _, configLocale := range *cfg.Store.Availabilities {
				if availablity.Name == configLocale {
					newAvailabilities = append(newAvailabilities, availablity)
				}
			}
		}

		ext.StoreAvailabilities = newAvailabilities
	}

	if cfg.Store.Categories != nil {
		newCategories := make([]accountApi.StoreCategory, 0)

		for _, category := range info.Categories {
			for _, configCategory := range *cfg.Store.Categories {
				if category.Name == configCategory {
					newCategories = append(newCategories, category)
				}
			}
		}

		ext.Categories = newCategories
	}

	if cfg.Store.Type != nil {
		for _, storeProductType := range info.ProductTypes {
			if storeProductType.Name == *cfg.Store.Type {
				ext.ProductType = storeProductType
			}
		}
	}

	if cfg.Store.AutomaticBugfixVersionCompatibility != nil {
		ext.AutomaticBugfixVersionCompatibility = *cfg.Store.AutomaticBugfixVersionCompatibility
	}

	for _, info := range ext.Infos {
		language := info.Locale.Name[0:2]

		if language == "de" {
			if cfg.Store.Info.German.Tags != nil {
				newTags := make([]accountApi.StoreTag, 0)

				for _, tag := range *cfg.Store.Info.German.Tags {
					newTags = append(newTags, accountApi.StoreTag{Name: tag})
				}

				info.Tags = newTags
			}

			if cfg.Store.Info.German.Videos != nil {
				newVideos := make([]accountApi.StoreVideo, 0)

				for _, video := range *cfg.Store.Info.German.Videos {
					newVideos = append(newVideos, accountApi.StoreVideo{URL: video})
				}

				info.Videos = newVideos
			}

			if cfg.Store.Info.German.Hightlight != nil {
				info.Highlights = strings.Join(*cfg.Store.Info.German.Hightlight, "\n")
			}

			if cfg.Store.Info.German.Features != nil {
				info.Features = strings.Join(*cfg.Store.Info.German.Features, "\n")
			}

			if cfg.Store.Info.German.Faq != nil {
				newFaq := make([]accountApi.StoreFaq, 0)

				for _, faq := range *cfg.Store.Info.German.Faq {
					newFaq = append(newFaq, accountApi.StoreFaq{Question: faq.Question, Answer: faq.Answer})
				}

				info.Faqs = newFaq
			}
		} else {
			if cfg.Store.Info.English.Tags != nil {
				newTags := make([]accountApi.StoreTag, 0)

				for _, tag := range *cfg.Store.Info.English.Tags {
					newTags = append(newTags, accountApi.StoreTag{Name: tag})
				}

				info.Tags = newTags
			}

			if cfg.Store.Info.English.Videos != nil {
				newVideos := make([]accountApi.StoreVideo, 0)

				for _, video := range *cfg.Store.Info.English.Videos {
					newVideos = append(newVideos, accountApi.StoreVideo{URL: video})
				}

				info.Videos = newVideos
			}

			if cfg.Store.Info.English.Hightlight != nil {
				info.Highlights = strings.Join(*cfg.Store.Info.English.Hightlight, "\n")
			}

			if cfg.Store.Info.English.Features != nil {
				info.Features = strings.Join(*cfg.Store.Info.English.Features, "\n")
			}

			if cfg.Store.Info.English.Faq != nil {
				newFaq := make([]accountApi.StoreFaq, 0)

				for _, faq := range *cfg.Store.Info.English.Faq {
					newFaq = append(newFaq, accountApi.StoreFaq{Question: faq.Question, Answer: faq.Answer})
				}

				info.Faqs = newFaq
			}
		}
	}
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionUpdateCmd)
}
