package redfish

import "time"

type OData struct {
	Id           *string `json:"@odata.id"`
	Type         *string `json:"@odata.type"`
	Context      *string `json:"@odata.context"`
	Members      []OData `json:"Members"`
	MembersCount int     `json:"Members@odata.count"`
}

type baseEndpoint struct {
	AccountService OData `json:"AccountService"`
	Chassis        OData `json:"Chassis"`
	Managers       OData `json:"Managers"`
	SessionService OData `json:"SessionService"`
	Systems        OData `json:"Systems"`
}

type sessionServiceEndpoint struct {
	Enabled        *bool  `json:"ServiceEnabled"`
	SessionTimeout int    `json:"SessionTimeout"`
	Sessions       *OData `json:"Sessions"`
}

type Status struct {
	State  string `json:"State"`
	Health string `json:"Health"`
}

type SystemProcessorSummary struct {
	Count  int    `json:"Count"`
	Status Status `json:"Status"`
}

type SystemData struct {
	UUID         *string `json:"UUID"`
	Status       Status  `json:"Status"`
	SerialNumber *string `json:"SerialNumber"`
	//ProcessorSummary  *SystemProcessorSummary  `json:"ProcessorSummary"`
	Processors *OData  `json:"Processors"`
	PowerState *string `json:"Powerstate"`
	Name       *string `json:"Name"`
	Model      *string `json:"Model"`
	//MemorySummary *SystemMemorySummary    `json:"MemorySummary"`
	Memory             *OData  `json:"Memory"`
	Manufacturer       *string `json:"Manufacturer"`
	LogServices        *OData  `json:"LogServices"`
	Id                 *string `json:"Id"`
	EthernetInterfaces *OData  `json:"EthernetInterfaces"`
	BIOSVersion        *string `json:"BiosVersion"`
	BIOS               *OData  `json:"Bios"`
	// Actions
	SelfEndpoint *string
}

type AccountService struct {
	AccountsEndpoint *OData `json:"Accounts"`
	RolesEndpoint    *OData `json:"Roles"`
}

type AccountData struct {
	Id       *string `json:"Id"`
	Name     *string `json:"Name"`
	UserName *string `json:"UserName"`
	Password *string `json:"Password"`
	RoleId   *string `json:"RoleId"`
	Enabled  *bool   `json:"Enabled"`
	Locked   *bool   `json:"Locked"`

	SelfEndpoint *string
}

type RoleData struct {
	Id                 *string  `json:"Id"`
	Name               *string  `json:"Name"`
	IsPredefined       *bool    `json:"IsPredefined"`
	Description        *string  `json:"Description"`
	AssignedPrivileges []string `json:"AssignedPrivileges"`
	//    OemPrivileges   []string    `json:"OemPrivileges"`
	SelfEndpoint *string
}

type ManagerData struct {
	Id              *string `json:"Id"`
	Name            *string `json:"Name"`
	ManagerType     *string `json:"ManagerType"`
	UUID            *string `json:"UUID"`
	Status          Status  `json:"Status"`
	FirmwareVersion *string `json:"FirmwareVersion"`

	/* futher data
	   Oem
	   VirtualMedia
	   SerialConsole
	   NetworkProtocol
	   GraphicalConsole
	   FirmwareVersion
	   EthernetInterfaces
	   Actions
	*/

	SelfEndpoint *string
}

const (
	REDFISH_GENERAL uint = iota
	REDFISH_HP
	REDFISH_HUAWEI
	REDFISH_INSPUR
	REDFISH_LENOVO
	REDFISH_SUPERMICRO
)

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
	MapSystensById() (map[string]*SystemData, error)
	MapSystemsByUuid() (map[string]*SystemData, error)
	GetAccounts() ([]string, error)
	GetAccountData(string) (*AccountData, error)
	MapAccountsByName() (map[string]*AccountData, error)
	MapAccountsById() (map[string]*AccountData, error)
	GetRoles() ([]string, error)
	GetRoleData(string) (*AccountData, error)
	MapRolesByName() (map[string]*RoleData, error)
	MapRolesById() (map[string]*RoleData, error)
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

	// Vendor "flavor"
	Flavor uint
}
