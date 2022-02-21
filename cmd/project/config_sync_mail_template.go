package project

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"shopware-cli/shop"
)

type MailTemplateSync struct{}

func (s MailTemplateSync) Push(ctx context.Context, client *shop.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	mailTemplates, err := fetchAllMailTemplates(ctx, client)
	if err != nil {
		return err
	}

	mailUpdates := make([]map[string]interface{}, 0)

	for _, external := range mailTemplates.Data {
		for _, configEntry := range config.Sync.MailTemplate {
			if external["id"] == configEntry.Id {
				mailUpdate := make(map[string]interface{})
				mailUpdate["id"] = configEntry.Id
				translationUpdates := make(map[string]map[string]interface{})

				translations := external["translations"].([]interface{})

				for _, translationRaw := range translations {
					translation := translationRaw.(map[string]interface{})

					if translation["language"] == nil {
						continue
					}

					language := translation["language"].(map[string]interface{})

					for _, configTranslation := range configEntry.Translations {
						if language["name"].(string) == configTranslation.Language {
							translationUpdate := make(map[string]interface{})

							if translation["senderName"].(string) != configTranslation.SenderName {
								translationUpdate["senderName"] = configTranslation.SenderName
							}

							if translation["subject"].(string) != configTranslation.Subject {
								translationUpdate["subject"] = configTranslation.Subject
							}

							if configTranslation.HTML != "" {
								if content, err := ioutil.ReadFile(configTranslation.HTML); err == nil {
									if translation["contentHtml"].(string) != string(content) {
										translationUpdate["contentHtml"] = string(content)
									}
								} else {
									log.Errorf("Cannot read file %s, with error: %s", configTranslation.HTML, err)
								}
							}

							if configTranslation.Plain != "" {
								if content, err := ioutil.ReadFile(configTranslation.Plain); err == nil {
									if translation["contentPlain"].(string) != string(content) {
										translationUpdate["contentPlain"] = string(content)
									}
								} else {
									log.Errorf("Cannot read file %s, with error: %s", configTranslation.Plain, err)
								}
							}

							localCustomFields, _ := json.Marshal(configTranslation.CustomFields)
							remoteCustomFields, _ := json.Marshal(translation["customFields"])

							if string(localCustomFields) != string(remoteCustomFields) {
								translationUpdate["customFields"] = configTranslation.CustomFields
							}

							if len(translationUpdate) > 0 {
								translationUpdates[translation["languageId"].(string)] = translationUpdate
							}
						}
					}

					if len(translationUpdates) > 0 {
						mailUpdate["translations"] = translationUpdates
					}
				}

				if len(mailUpdate) > 1 {
					mailUpdates = append(mailUpdates, mailUpdate)
				}
			}
		}
	}

	if len(mailUpdates) > 0 {
		operation.Operations["update-mail-template"] = shop.SyncOperation{
			Action:  "upsert",
			Entity:  "mail_template",
			Payload: mailUpdates,
		}
	}

	return nil
}

func (s MailTemplateSync) Pull(ctx context.Context, client *shop.Client, config *shop.Config) error {
	mailTemplates, err := fetchAllMailTemplates(ctx, client)
	if err != nil {
		return err
	}

	config.Sync.MailTemplate = make([]shop.MailTemplate, 0)

	for _, row := range mailTemplates.Data {
		if row["mailTemplateType"] == nil {
			log.Infof("mail_template entity with id %s does not have a type. Skipping", row["id"])
			continue
		}

		mailType := row["mailTemplateType"].(map[string]interface{})
		translations := row["translations"].([]interface{})

		cfg := shop.MailTemplate{
			Id:           row["id"].(string),
			Translations: []shop.MailTemplateTranslation{},
		}

		for _, translationRaw := range translations {
			translation := translationRaw.(map[string]interface{})

			if translation["language"] == nil {
				continue
			}

			language := translation["language"].(map[string]interface{})
			configKey := language["name"].(string)

			htmLFilePath := fmt.Sprintf(".shopware-cli/mail-template/%s/%s-html.twig", mailType["technicalName"].(string), configKey)
			plainFilePath := fmt.Sprintf(".shopware-cli/mail-template/%s/%s-plain.twig", mailType["technicalName"].(string), configKey)
			dir := filepath.Dir(htmLFilePath)

			if _, err := os.Stat(dir); os.IsNotExist(err) {
				if createErr := os.MkdirAll(dir, os.ModePerm); createErr != nil {
					return createErr
				}
			}

			cfgLang := shop.MailTemplateTranslation{
				Language:     configKey,
				SenderName:   translation["senderName"].(string),
				Subject:      translation["subject"].(string),
				HTML:         htmLFilePath,
				Plain:        plainFilePath,
				CustomFields: translation["customFields"],
			}

			if err := ioutil.WriteFile(htmLFilePath, []byte(translation["contentHtml"].(string)), os.ModePerm); err != nil {
				return err
			}

			if err := ioutil.WriteFile(plainFilePath, []byte(translation["contentPlain"].(string)), os.ModePerm); err != nil {
				return err
			}

			cfg.Translations = append(cfg.Translations, cfgLang)
		}

		config.Sync.MailTemplate = append(config.Sync.MailTemplate, cfg)
	}

	return nil
}

func fetchAllMailTemplates(ctx context.Context, client *shop.Client) (*shop.SearchResponse, error) {
	criteria := shop.Criteria{}
	criteria.Includes = map[string][]string{
		"mail_template":             {"id", "mailTemplateType", "translations"},
		"mail_template_type":        {"technicalName"},
		"mail_template_translation": {"senderName", "subject", "contentHtml", "contentPlain", "language", "languageId"},
		"language":                  {"name"},
	}
	criteria.Associations = map[string]shop.Criteria{"mailTemplateType": {}, "translations": {Associations: map[string]shop.Criteria{"language": {}}}}

	return client.SearchAll(ctx, "mail_template", criteria)
}
