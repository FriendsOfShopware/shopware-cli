package account

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	goldmarkExtension "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	accountApi "github.com/FriendsOfShopware/shopware-cli/account-api"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var accountCompanyProducerExtensionInfoPushCmd = &cobra.Command{
	Use:   "push [zip or path]",
	Short: "Update store information of extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}

		stat, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot open file: %w", err)
		}

		var zipExt extension.Extension

		if stat.IsDir() {
			zipExt, err = extension.GetExtensionByFolder(path)
		} else {
			zipExt, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			return fmt.Errorf("cannot open extension: %w", err)
		}

		zipName, err := zipExt.GetName()
		if err != nil {
			return fmt.Errorf("cannot get name: %w", err)
		}

		p, err := services.AccountClient.Producer(cmd.Context())

		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		storeExt, err := p.GetExtensionByName(cmd.Context(), zipName)

		if err != nil {
			return fmt.Errorf("cannot get store extension: %w", err)
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

		info, err := p.GetExtensionGeneralInfo(cmd.Context())
		if err != nil {
			return fmt.Errorf("cannot get general info: %w", err)
		}

		extCfg, err := extension.ReadExtensionConfig(zipExt.GetPath())
		if err != nil {
			return fmt.Errorf("cannot read extension config: %w", err)
		}

		if extCfg != nil {
			if extCfg.Store.Icon != nil {
				err := p.UpdateExtensionIcon(cmd.Context(), storeExt.Id, fmt.Sprintf("%s/%s", zipExt.GetPath(), *extCfg.Store.Icon))
				if err != nil {
					return fmt.Errorf("cannot update extension icon due error: %w", err)
				}
			}

			if extCfg.Store.Images != nil {
				images, err := p.GetExtensionImages(cmd.Context(), storeExt.Id)
				if err != nil {
					return fmt.Errorf("cannot get images from remote server: %w", err)
				}

				for _, image := range images {
					err := p.DeleteExtensionImages(cmd.Context(), storeExt.Id, image.Id)

					if err != nil {
						return fmt.Errorf("cannot extension image: %w", err)
					}
				}

				for _, configImage := range *extCfg.Store.Images {
					apiImage, err := p.AddExtensionImage(cmd.Context(), storeExt.Id, fmt.Sprintf("%s/%s", zipExt.GetPath(), configImage.File))

					if err != nil {
						return fmt.Errorf("cannot upload image to extension: %w", err)
					}

					apiImage.Priority = configImage.Priority
					apiImage.Details[0].Activated = configImage.Activate.German
					apiImage.Details[0].Preview = configImage.Preview.German

					apiImage.Details[1].Activated = configImage.Activate.English
					apiImage.Details[1].Preview = configImage.Preview.English

					err = p.UpdateExtensionImage(cmd.Context(), storeExt.Id, apiImage)

					if err != nil {
						return fmt.Errorf("cannot update image information of extension: %w", err)
					}
				}
			}

			if err := updateStoreInfo(storeExt, zipExt, extCfg, info); err != nil {
				return fmt.Errorf("cannot update store information: %w", err)
			}
		}

		err = p.UpdateExtension(cmd.Context(), storeExt)

		if err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Store information has been updated")

		return nil
	},
}

func updateStoreInfo(ext *accountApi.Extension, zipExt extension.Extension, cfg *extension.Config, info *accountApi.ExtensionGeneralInformation) error { //nolint:gocyclo
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

		for _, availability := range info.StoreAvailabilities {
			for _, configLocale := range *cfg.Store.Availabilities {
				if availability.Name == configLocale {
					newAvailabilities = append(newAvailabilities, availability)
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
		for i, storeProductType := range info.ProductTypes {
			if storeProductType.Name == *cfg.Store.Type {
				ext.ProductType = &info.ProductTypes[i]
			}
		}
	}

	if cfg.Store.AutomaticBugfixVersionCompatibility != nil {
		ext.AutomaticBugfixVersionCompatibility = *cfg.Store.AutomaticBugfixVersionCompatibility
	}

	for _, info := range ext.Infos {
		language := info.Locale.Name[0:2]

		storeTags := getTranslation(language, cfg.Store.Tags)
		if storeTags != nil {
			var newTags []accountApi.StoreTag
			for _, tag := range *storeTags {
				newTags = append(newTags, accountApi.StoreTag{Name: tag})
			}

			info.Tags = newTags
		}

		storeVideos := getTranslation(language, cfg.Store.Videos)
		if storeVideos != nil {
			var newVideos []accountApi.StoreVideo
			for _, video := range *storeVideos {
				newVideos = append(newVideos, accountApi.StoreVideo{URL: video})
			}

			info.Videos = newVideos
		}

		storeHighlights := getTranslation(language, cfg.Store.Highlights)
		if storeHighlights != nil {
			info.Highlights = strings.Join(*storeHighlights, "\n")
		}

		storeFeatures := getTranslation(language, cfg.Store.Features)
		if storeFeatures != nil {
			info.Features = strings.Join(*storeFeatures, "\n")
		}

		storeFaqs := getTranslation(language, cfg.Store.Faq)
		if storeFaqs != nil {
			var newFaq []accountApi.StoreFaq
			for _, faq := range *storeFaqs {
				newFaq = append(newFaq, accountApi.StoreFaq{Question: faq.Question, Answer: faq.Answer})
			}

			info.Faqs = newFaq
		}

		var err error

		storeDescription := getTranslation(language, cfg.Store.Description)
		if storeDescription != nil {
			info.Description, err = parseInlineablePath(*storeDescription, zipExt.GetPath())

			if err != nil {
				return err
			}
		}

		storeManual := getTranslation(language, cfg.Store.InstallationManual)
		if storeManual != nil {
			info.InstallationManual, err = parseInlineablePath(*storeManual, zipExt.GetPath())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getTranslation[T extension.Translatable](language string, config extension.ConfigTranslated[T]) *T {
	if language == "de" {
		return config.German
	} else if language == "en" {
		return config.English
	}

	return nil
}

func init() {
	accountCompanyProducerExtensionInfoCmd.AddCommand(accountCompanyProducerExtensionInfoPushCmd)
}

func parseInlineablePath(path, extensionDir string) (string, error) {
	if !strings.HasPrefix(path, "file:") {
		return path, nil
	}

	filePath := fmt.Sprintf("%s/%s", extensionDir, strings.TrimPrefix(path, "file:"))

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error reading file at path %s with error: %v", filePath, err)
	}

	if filepath.Ext(filePath) != ".md" {
		return string(content), nil
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
		return "", fmt.Errorf("Cannot convert file at path %s from markdown to html with error: %v", filePath, err)
	}

	return buf.String(), nil
}
