package telemetry

import (
	"os"

	"github.com/posthog/posthog-go"
)

var (
	client     posthog.Client
	distinctId = ""
)

func Init() {
	doNotTrack := os.Getenv("DO_NOT_TRACK")

	if doNotTrack == "1" {
		return
	}

	var err error
	client, err = posthog.NewWithConfig(
		"phc_bkTsf3bZw70kXFoqMj3Qkeo3dtzC3x1JvuhamdJI8mJ",
		posthog.Config{
			Endpoint: "https://eu.posthog.com",
		},
	)
	if err != nil {
		return
	}

	distinctId = gatherDistinctId()
}

func Track(event string, properties map[string]interface{}) {
	if client == nil {
		return
	}

	client.Enqueue(posthog.Capture{
		Event:      event,
		DistinctId: distinctId,
		Properties: properties,
	})
}

func Close() {
	if client == nil {
		return
	}

	client.Close()
}
