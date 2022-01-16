package cmd

import (
	"bytes"
	"fmt"
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	goldmarkExtension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	accountApi "shopware-cli/account-api"
	"shopware-cli/extension"
	"strings"
)

var accountCompanyProducerExtensionInfoPushCmd = &cobra.Command{
	Use:   "push [zip or path]",
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
			if extCfg.Store.Icon != nil {
				err := p.UpdateExtensionIcon(storeExt.Id, fmt.Sprintf("%s/%s", zipExt.GetPath(), *extCfg.Store.Icon))
				if err != nil {
					log.Fatalln(fmt.Errorf("cannot update extension icon due error: %v", err))
				}
			}

			if extCfg.Store.Images != nil {
				images, err := p.GetExtensionImages(storeExt.Id)
				if err != nil {
					log.Fatalln(fmt.Errorf("cannot get images from remote server: %v", err))
				}

				for _, image := range images {
					err := p.DeleteExtensionImages(storeExt.Id, image.Id)

					if err != nil {
						log.Fatalln(fmt.Errorf("cannot extension image: %v", err))
					}
				}

				for _, configImage := range *extCfg.Store.Images {
					apiImage, err := p.AddExtensionImage(storeExt.Id, fmt.Sprintf("%s/%s", zipExt.GetPath(), configImage.File))

					if err != nil {
						log.Fatalln(fmt.Errorf("cannot upload image to extension: %v", err))
					}

					apiImage.Priority = configImage.Priority
					apiImage.Details[0].Activated = configImage.Activate.German
					apiImage.Details[0].Preview = configImage.Preview.German

					apiImage.Details[1].Activated = configImage.Activate.English
					apiImage.Details[1].Preview = configImage.Preview.English

					err = p.UpdateExtensionImage(storeExt.Id, apiImage)

					if err != nil {
						log.Fatalln(fmt.Errorf("cannot update image information of extension: %v", err))
					}
				}
			}

			updateStoreInfo(storeExt, zipExt, extCfg, info)
		}

		err = p.UpdateExtension(storeExt)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		termColor.Green("Store information has been updated")
	},
}

func updateStoreInfo(ext *accountApi.Extension, zipExt extension.Extension, cfg *extension.Config, info *accountApi.ExtensionGeneralInformation) {
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
			if cfg.Store.Tags.German != nil {
				newTags := make([]accountApi.StoreTag, 0)

				for _, tag := range *cfg.Store.Tags.German {
					newTags = append(newTags, accountApi.StoreTag{Name: tag})
				}

				info.Tags = newTags
			}

			if cfg.Store.Videos.German != nil {
				newVideos := make([]accountApi.StoreVideo, 0)

				for _, video := range *cfg.Store.Videos.German {
					newVideos = append(newVideos, accountApi.StoreVideo{URL: video})
				}

				info.Videos = newVideos
			}

			if cfg.Store.Highlights.German != nil {
				info.Highlights = strings.Join(*cfg.Store.Highlights.German, "\n")
			}

			if cfg.Store.Features.German != nil {
				info.Features = strings.Join(*cfg.Store.Features.German, "\n")
			}

			if cfg.Store.Faq.German != nil {
				newFaq := make([]accountApi.StoreFaq, 0)

				for _, faq := range *cfg.Store.Faq.German {
					newFaq = append(newFaq, accountApi.StoreFaq{Question: faq.Question, Answer: faq.Answer})
				}

				info.Faqs = newFaq
			}

			if cfg.Store.Description.German != nil {
				info.Description = parseInlineablePath(*cfg.Store.Description.German, zipExt.GetPath())
			}

			if cfg.Store.InstallationManual.German != nil {
				info.InstallationManual = parseInlineablePath(*cfg.Store.InstallationManual.German, zipExt.GetPath())
			}
		} else {
			if cfg.Store.Tags.English != nil {
				newTags := make([]accountApi.StoreTag, 0)

				for _, tag := range *cfg.Store.Tags.English {
					newTags = append(newTags, accountApi.StoreTag{Name: tag})
				}

				info.Tags = newTags
			}

			if cfg.Store.Videos.English != nil {
				newVideos := make([]accountApi.StoreVideo, 0)

				for _, video := range *cfg.Store.Videos.English {
					newVideos = append(newVideos, accountApi.StoreVideo{URL: video})
				}

				info.Videos = newVideos
			}

			if cfg.Store.Highlights.English != nil {
				info.Highlights = strings.Join(*cfg.Store.Highlights.English, "\n")
			}

			if cfg.Store.Features.English != nil {
				info.Features = strings.Join(*cfg.Store.Features.English, "\n")
			}

			if cfg.Store.Faq.English != nil {
				newFaq := make([]accountApi.StoreFaq, 0)

				for _, faq := range *cfg.Store.Faq.English {
					newFaq = append(newFaq, accountApi.StoreFaq{Question: faq.Question, Answer: faq.Answer})
				}

				info.Faqs = newFaq
			}

			if cfg.Store.Description.English != nil {
				info.Description = parseInlineablePath(*cfg.Store.Description.English, zipExt.GetPath())
			}

			if cfg.Store.InstallationManual.English != nil {
				info.InstallationManual = parseInlineablePath(*cfg.Store.InstallationManual.English, zipExt.GetPath())
			}
		}
	}
}

func init() {
	accountCompanyProducerExtensionInfoCmd.AddCommand(accountCompanyProducerExtensionInfoPushCmd)
}

func parseInlineablePath(path, extensionDir string) string {
	if !strings.HasPrefix(path, "file:") {
		return path
	}

	filePath := fmt.Sprintf("%s/%s", extensionDir, strings.TrimPrefix(path, "file:"))

	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Error reading file at path %s with error: %v", filePath, err))
	}

	if filepath.Ext(filePath) != ".md" {
		return string(content)
	}

	md := goldmark.New(
		goldmark.WithExtensions(goldmarkExtension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	err = md.Convert(content, &buf)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot convert file at path %s from markdown to html with error: %v", filePath, err))
	}

	return buf.String()
}
