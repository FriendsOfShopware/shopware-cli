package project

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/shop"
)

var projectConfigPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Synchronizes your local config to the external shop",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		autoApprove, _ := cmd.PersistentFlags().GetBool("auto-approve")

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		operation := &ConfigSyncOperation{
			Upsert: map[string][]map[string]interface{}{"system_config": {}},
		}

		for _, applyer := range NewSyncApplyers() {
			if err := applyer.Push(cmd.Context(), client, cfg, operation); err != nil {
				return err
			}
		}

		if !operation.HasChanges() {
			log.Infof("Configuration is up to date")
			return nil
		}

		if len(operation.Upsert) > 0 {
			log.Println("Following entities will be written")

			for entity, values := range operation.Upsert {
				log.Printf("Entity: %s", entity)

				content, _ := json.Marshal(values)

				log.Printf("Payload: %s", string(content))
			}
		}

		if !autoApprove {
			p := promptui.Prompt{
				Label:     "You want to apply these changes to your Shop?",
				IsConfirm: true,
			}

			if _, err := p.Run(); err != nil {
				return err
			}
		}

		syncApiOperation := make(map[string]shop.SyncOperation)

		for entity, values := range operation.Upsert {
			syncApiOperation[fmt.Sprintf("upsert-%s", entity)] = shop.SyncOperation{
				Entity:  entity,
				Action:  "upsert",
				Payload: values,
			}
		}

		for entity, values := range operation.Delete {
			syncApiOperation[fmt.Sprintf("delete-%s", entity)] = shop.SyncOperation{
				Entity:  entity,
				Action:  "delete",
				Payload: values,
			}
		}

		if err := client.Sync(cmd.Context(), syncApiOperation); err != nil {
			return err
		}

		log.Infof("Configuration has been applied to remote")

		return nil
	},
}

func init() {
	projectConfigCmd.AddCommand(projectConfigPushCmd)
	projectConfigPushCmd.PersistentFlags().Bool("auto-approve", false, "Skips the confirmation")
}
