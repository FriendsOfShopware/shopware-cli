package main

import (
	"context"

	"github.com/FriendsOfShopware/shopware-cli/cmd"
	"github.com/FriendsOfShopware/shopware-cli/internal/telemetry"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

func main() {
	telemetry.Init()
	defer telemetry.Close()

	logger := logging.NewLogger()
	cmd.Execute(logging.WithLogger(context.Background(), logger))
}
