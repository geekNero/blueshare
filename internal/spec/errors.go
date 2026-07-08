package spec

type Error int

type ErrorResponse struct {
	Error
	Err error
}

const (
	// Common errors
	ErrFailedToEnableAdapter = 1001
	ErrUnknown               = 1002

	// Scanning errors
	ErrScanAlreadyInProgress = 2000

	// Broadcasting errors
	ErrInvalidMessage = 3000
)

var (
	ErrorMap = map[Error]string{
		ErrScanAlreadyInProgress: "scan is already in-progress",
		ErrFailedToEnableAdapter: "failed to enable adapter, is bluetooth on?",
		ErrInvalidMessage:        "input message contains invalid characters",
		ErrUnknown:               "unknown error",
	}
)
