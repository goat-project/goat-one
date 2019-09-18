package constants

// constants with error messages
const (
	ErrCreatePrepReaderNil  = "error create Preparer when reader is nil"
	ErrCreatePrepLimiterNil = "error create Preparer when limiter is nil"
	ErrCreatePrepConnNil    = "error create Preparer when gRPC client connection is nil"

	ErrPrepEmptyNetUser = "error prepare empty NetUser"
	ErrPrepNoNetUser    = "error get id, unable to prepare network record"

	ErrPrepIPv4 = "unable to prepare ipv4 network record"
	ErrPrepIPv6 = "unable to prepare ipv6 network record"

	ErrPrepWrite = "error write network record"

	ErrNoSiteName  = "no site name in configuration"
	ErrNoCloudType = "no cloud type in configuration"
	ErrNoGroupName = "no group name"

	ErrCreateProcReaderNil = "error create Processor when Reader is nil"
)
