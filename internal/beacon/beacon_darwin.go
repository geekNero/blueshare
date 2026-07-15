//go:build darwin

package beacon

import (
	"blueshare/internal/spec"
)

func Broadcast(c *BeaconCmd) spec.ErrorResponse {
	return spec.ErrorResponse{
		Error: spec.ErrNotImplemented,
	}
}

func Stop() {
}

func FetchError() error {
	return nil
}
