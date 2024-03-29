package main

import (
	"context"

	"github.com/FriendsOfShopware/shopware-cli/cmd"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

func main() {
	logger := logging.NewLogger()
	cmd.Execute(logging.WithLogger(context.Background(), logger))
}
