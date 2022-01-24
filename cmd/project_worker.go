package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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

		cancelCtx, cancel := context.WithCancel(ctx)
		cancelOnTermination(cancel)

		var wg sync.WaitGroup
		for a := 0; a < workerAmount; a++ {
			wg.Add(1)
			go func() {
				for {
					cmd := exec.CommandContext(cancelCtx, "php", "bin/console", "messenger:consume", "-vvv")
					cmd.Dir = projectRoot
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}
				}
			}()
		}

		wg.Wait()

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectWorkerCmd)
}

func cancelOnTermination(cancel context.CancelFunc) {
	log.Println("setting up a signal handler")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		log.Printf("received SIGTERM %v\n", <-s)
		cancel()
	}()
}
