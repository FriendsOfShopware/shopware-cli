package project

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/spf13/cobra"
)

var projectWorkerCmd = &cobra.Command{
	Use:   "worker [amount]",
	Short: "Runs the Symfony Worker in Background",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		var projectRoot string
		var err error
		workerAmount := 1

		if projectRoot, err = findClosestShopwareProject(); err != nil {
			return err
		}

		if len(args) > 0 {
			workerAmount, err = strconv.Atoi(args[0])

			if err != nil {
				return err
			}
		}

		cancelCtx, cancel := context.WithCancel(cobraCmd.Context())
		cancelOnTermination(cancelCtx, cancel)

		var wg sync.WaitGroup
		for a := 0; a < workerAmount; a++ {
			wg.Add(1)
			go func(ctx context.Context) {
				for {
					cmd := exec.CommandContext(cancelCtx, "php", "bin/console", "messenger:consume", "-vvv")
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
