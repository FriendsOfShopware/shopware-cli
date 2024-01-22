package project

import "github.com/spf13/cobra"

var dockerRootCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker Tools",
}

func init() {
	projectRootCmd.AddCommand(dockerRootCmd)
}
