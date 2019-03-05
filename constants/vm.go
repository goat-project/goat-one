package constants

// prefix for virtual machine subcommands
const cfgVMPrefix = "vm."

// constants for virtual machine subcommand
const (
	// CfgSiteName represents string of virtual machine site name
	CfgSiteName = cfgVMPrefix + "site-name"
	// CfgCloudType represents string of virtual machine cloud type
	CfgCloudType = cfgVMPrefix + "cloud-type"
	// CfgCloudComputeService represents string of virtual machine cloud compute service
	CfgCloudComputeService = cfgVMPrefix + "cloud-compute-service"
)
