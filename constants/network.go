package constants

// prefix for network subcommands
const cfgNetworkPrefix = "network."

// constants for network subcommand
const (
	// CfgNetworkSiteName represents string of network site name
	CfgNetworkSiteName = cfgNetworkPrefix + "site-name"
	// CfgNetworkCloudType represents string of network cloud type
	CfgNetworkCloudType = cfgNetworkPrefix + "cloud-type"
	// CfgNetworkCloudComputeService represents string of network cloud compute service
	CfgNetworkCloudComputeService = cfgNetworkPrefix + "cloud-compute-service"
)

// constants for network errors
const (
	ErrCreatePrepLimiterNil = "error create Preparer when limiter is nil"
	ErrCreatePrepConnNil    = "error create Preparer when gRPC client connection is nil"
	ErrPrepEmptyNetUser     = "error prepare empty NetUser"
	ErrPrepNoNetUser        = "error get id, unable to prepare network record"
	ErrPrepIPv4             = "unable to prepare ipv4 network record"
	ErrPrepIPv6             = "unable to prepare ipv6 network record"
	ErrPrepWrite            = "error write network record"
	ErrNoSiteName           = "no site name in configuration"
	ErrNoCloudType          = "no cloud type in configuration"
	ErrNoGroupName          = "no group name"
)
