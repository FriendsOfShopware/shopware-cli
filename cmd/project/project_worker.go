package project

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/FriendsOfShopware/shopware-cli/shop"

	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var projectWorkerCmd = &cobra.Command{
	Use:   "worker [amount]",
	Short: "Runs the Symfony Worker in Background",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		var projectRoot string
		var err error
		workerAmount := 1

		isVerbose, _ := cobraCmd.Flags().GetBool("verbose")
		queuesToConsume, _ := cobraCmd.Flags().GetString("queue")
		memoryLimit, _ := cobraCmd.Flags().GetString("memory-limit")

		if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		if len(args) > 0 {
			workerAmount, err = strconv.Atoi(args[0])

			if err != nil {
				return err
			}
		}

		if memoryLimit == "" {
			memoryLimit = "512M"
		}

		cancelCtx, cancel := context.WithCancel(cobraCmd.Context())
		cancelOnTermination(cancelCtx, cancel)

		consumeArgs := []string{"bin/console", "messenger:consume", fmt.Sprintf("--memory-limit=%s", memoryLimit)}

		if queuesToConsume == "" {
			if is, _ := shop.IsShopwareVersion(projectRoot, ">=6.5"); is {
				consumeArgs = append(consumeArgs, "async", "failed")
			}
		} else {
			consumeArgs = append(consumeArgs, strings.Split(queuesToConsume, ",")...)
		}

		if isVerbose {
			consumeArgs = append(consumeArgs, "-vvv")
		}

		var wg sync.WaitGroup
		for a := 0; a < workerAmount; a++ {
			wg.Add(1)
			go func(ctx context.Context) {
				for {
					cmd := exec.CommandContext(cancelCtx, "php", consumeArgs...)
					cmd.Dir = projectRoot
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						logging.FromContext(ctx).Fatal(err)
					}
				}
			}(cancelCtx)
		}

		wg.Wait()

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectWorkerCmd)
	projectWorkerCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	projectWorkerCmd.PersistentFlags().String("queue", "", "Queues to consume")
	projectWorkerCmd.PersistentFlags().String("memory-limit", "", "Memory Limit")
}

func cancelOnTermination(ctx context.Context, cancel context.CancelFunc) {
	logging.FromContext(ctx).Infof("setting up a signal handler")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		logging.FromContext(ctx).Infof("received SIGTERM %v\n", <-s)
		cancel()
	}()
}
