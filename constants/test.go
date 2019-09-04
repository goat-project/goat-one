package constants

// Constants for testing
const (
	// OpenNebulaEndpoint is endpoint to test OpenNebula
	OpenNebulaEndpoint = "http://192.168.122.208:2633/RPC2"
	// WrongOpenNebulaEndpoint is wrong endpoint to test OpenNebula
	WrongOpenNebulaEndpoint = "http://192.168.122.111:2633/RPC2"

	name          = "oneadmin"
	password      = "opennebula"
	wrongName     = "admin"
	wrongPassword = "nebula"

	// Token to connect to OpenNebula test account
	Token = name + ":" + password
	// WrongNameToken to connect to OpenNebula test account with wrong username
	WrongNameToken = wrongName + ":" + password
	// WrongPswdToken to connect to OpenNebula test account with wrong password
	WrongPswdToken = name + ":" + wrongPassword

	// NumTestedNetworks is a number of tested networks
	NumTestedNetworks = 2

	// BigPageOffset is a number bigger than number of networks divided by page size
	BigPageOffset = 1000
	// NegPageOffset is a negative page offset
	NegPageOffset = -10
)
