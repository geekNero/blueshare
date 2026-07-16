package spec

import "time"

var (
	SocketAddr = "blueshare.sock"

	// Suffix of this string has to be a 7-digit number
	CustomUUIDString = "12345678-1234-5678-1234-56789%s"
	// CustomUUID       bluetooth.UUID

	// Limits
	MaximumMessageLength                   = 1024
	BroadcastWindow                        = 120 // seconds
	DefaultPerPacketDuration time.Duration = 20  // seconds
	ByteSizeManufacturerData               = 27
	ByteSizeServiceData                    = 12
)
