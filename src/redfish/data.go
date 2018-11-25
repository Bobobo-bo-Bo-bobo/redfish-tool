package redfish

import "time"

type RedfishConfiguration struct {
	Hostname        string
	Port            int
	Username        string
	Password        string
	AuthToken       *string
	SessionLocation *string
	Timeout         time.Duration
	InsecureSSL     bool
	Verbose         bool
	RawBaseContent  string

	// endpoints
	AccountService string
	Chassis        string
	Managers       string
	SessionService string
	Sessions       string
	Systems        string
}

type Result struct {
	RawContent *string
	RawHeaders *string
}

type BaseRedfish interface {
	Initialize(*RedfishConfiguration) error
	Login(*RedfishConfiguration) error
	Logout(*RedfishConfiguration) error
	GetSystems(*RedfishConfiguration) ([]string, error)
	GetSystemData(*RedfishConfiguration, string) (*SystemData, error)
	GetAccounts(*RedfishConfiguration) ([]string, error)
	GetAccountData(*RedfishConfiguration, string) (*AccountData, error)
	GetRoles(*RedfishConfiguration) ([]string, error)
	GetRoleData(*RedfishConfiguration, string) (*AccountData, error)
}

type Redfish struct {
}
