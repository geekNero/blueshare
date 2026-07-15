package spec

type Error int

type ErrorResponse struct {
	Error
	ErrMsg string
}

const (
	// Common errors
	ErrFailedToEnableAdapter = 1001
	ErrUnknown               = 1002

	// Scanning errors
	ErrScanAlreadyInProgress = 2000

	// Broadcasting errors
	ErrInvalidMessage             = 3000
	ErrBroadcastAlreadyInProgress = 3001
	ErrNotImplemented             = 3002
)

var ErrorMap = map[Error]string{
	ErrScanAlreadyInProgress:      "scan is already in-progress",
	ErrFailedToEnableAdapter:      "failed to enable adapter, is bluetooth on?",
	ErrInvalidMessage:             "input message contains invalid characters",
	ErrUnknown:                    "unknown error",
	ErrNotImplemented:             "broadcast is not implemented for this operating system",
	ErrBroadcastAlreadyInProgress: "a message is already under broadcast, stop the existing broadcast or force a new one",
}
