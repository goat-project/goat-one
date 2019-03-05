package constants

// global constants
const (
	// CfgIdentifier represents string identifier of a goat-one instance
	CfgIdentifier = "identifier"
	// CfgRecordsFrom represents time which records are filtered from
	CfgRecordsFrom = "records-from"
	// CfgRecordsTo represents time which records are filtered to
	CfgRecordsTo = "records-to"
	// CfgRecordsFrom represents duration which records are filtered for
	CfgRecordsForPeriod = "records-for-period"
	// CfgEndpoint represents string (address:port) of goat server endpoint
	CfgEndpoint = "endpoint"
	// CfgOpennebulaEndpoint represents string (address:port) of OpenNebula endpoint
	CfgOpennebulaEndpoint = "opennebula-endpoint"
	// CfgOpennebulaSecret represents string (username:password) of user login
	CfgOpennebulaSecret = "opennebula-secret" // nolint: gosec
	// CfgOpennebulaTimeout represents duration (timeout) for OpenNebula calls
	CfgOpennebulaTimeout = "opennebula-timeout"
	// CfgDebug represents true for debug mode; false otherwise
	CfgDebug = "debug"
)
