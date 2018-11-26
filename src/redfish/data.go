package redfish

import "time"

type Result struct {
	RawContent *string
	RawHeaders *string
}

type BaseRedfish interface {
	Initialize() error
	Login() error
	Logout() error
	GetSystems() ([]string, error)
	GetSystemData(string) (*SystemData, error)
	GetAccounts() ([]string, error)
	GetAccountData(string) (*AccountData, error)
	MapAccountsByName() (map[string]*AccountData, error)
	MapAccountsById() (map[string]*AccountData, error)
	GetRoles() ([]string, error)
	GetRoleData(string) (*AccountData, error)
}

type Redfish struct {
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
