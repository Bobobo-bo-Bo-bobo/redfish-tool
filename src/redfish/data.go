package redfish

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const UserAgent string = "go-redfish/1.0.0"

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
	State        *string `json:"State"`
	Health       *string `json:"Health"`
	HealthRollUp *string `json:"HealthRollUp"`
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

type ManagerActionsData struct {
	ManagerReset LinkTargets `json:"#Manager.Reset"`
}

type ManagerData struct {
	Id              *string         `json:"Id"`
	Name            *string         `json:"Name"`
	ManagerType     *string         `json:"ManagerType"`
	UUID            *string         `json:"UUID"`
	Status          Status          `json:"Status"`
	FirmwareVersion *string         `json:"FirmwareVersion"`
	Oem             json.RawMessage `json:"Oem"`
	Actions         json.RawMessage `json:"Actions"` // may contain vendor specific data and endpoints

	/* futher data
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

type X509CertInfo struct {
	Issuer         *string `json:"Issuer"`
	SerialNumber   *string `json:"SerialNumber"`
	Subject        *string `json:"Subject"`
	ValidNotAfter  *string `json:"ValidNotAfter"`
	ValidNotBefore *string `json:"ValidNotBefore"`
}

type LinkTargets struct {
	Target     *string `json:"target"`
	ActionInfo *string `json:"@Redfish.ActionInfo"`
}

// HP/HPE: Oem data for Manager endpoint and SecurityService endpoint
type ManagerDataOemHpLicense struct {
	Key    *string `json:"LicenseKey"`
	String *string `json:"LicenseString"`
	Type   *string `json:"LicenseType"`
}

type ManagerDataOemHpFederationConfig struct {
	IPv6MulticastScope            *string `json:"IPv6MulticastScope"`
	MulticastAnnouncementInterval *int64  `json:"MulticastAnnouncementInterval"`
	MulticastDiscovery            *string `json:"MulticastDiscovery"`
	MulticastTimeToLive           *int64  `json:"MulticastTimeToLive"`
	ILOFederationManagement       *string `json:"iLOFederationManagement"`
}

type ManagerDataOemHpFirmwareData struct {
	Date         *string `json:"Date"`
	DebugBuild   *bool   `json:"DebugBuild"`
	MajorVersion *uint64 `json:"MajorVersion"`
	MinorVersion *uint64 `json:"MinorVersion"`
	Time         *string `json:"Time"`
	Version      *string `json:"Version"`
}

type ManagerDataOemHpFirmware struct {
	Current ManagerDataOemHpFirmwareData `json:"Current"`
}

// NOTE: The result for HP/HPE are different depending if the HTTP header
// OData-Version is set or not. If OData-Version is _NOT_ set data are returned in
// .Oem.Hp.links with endpoints in "href". If OData-Version is set
// data are returned in .Oem.Hp.Links (note different case!) and endpoints are
// defined as @odata.id. We always set "OData-Version: 4.0"
type ManagerDataOemHpLinks struct {
	ActiveHealthSystem   OData `json:"ActiveHealthSystem"`
	DateTimeService      OData `json:"DateTimeService"`
	EmbeddedMediaService OData `json:"EmbeddedMediaService"`
	FederationDispatch   OData `json:"FederationDispatch"`
	FederationGroups     OData `json:"FederationGroups"`
	FederationPeers      OData `json:"FederationPeers"`
	LicenseService       OData `json:"LicenseService"`
	SecurityService      OData `json:"SecurityService"`
	UpdateService        OData `json:"UpdateService"`
	VSPLogLocation       OData `json:"VSPLogLocation"`
}

type _managerDataOemHp struct {
	FederationConfig ManagerDataOemHpFederationConfig `json:"FederationConfig"`
	Firmware         ManagerDataOemHpFirmware         `json:"Firmware"`
	License          ManagerDataOemHpLicense          `json:"License"`
	Type             *string                          `json:"Type"`
	Links            ManagerDataOemHpLinks            `json:"Links"`
}

type ManagerDataOemHp struct {
	Hp _managerDataOemHp `json:"Hp"`
}

type SecurityServiceDataOemHpLinks struct {
	ESKM      OData `json:"ESKM"`
	HttpsCert OData `json:"HttpsCert"`
	SSO       OData `json:"SSO"`
}

type SecurityServiceDataOemHp struct {
	Id    *string                       `json:"Id"`
	Type  *string                       `json:"Type"`
	Links SecurityServiceDataOemHpLinks `json:"Links"`
}

type HttpsCertActionsOemHp struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HpHttpsCert.ImportCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

type HttpsCertDataOemHp struct {
	CSR     *string               `json:"CertificateSigningRequest"`
	Id      *string               `json:"Id"`
	Actions HttpsCertActionsOemHp `json:"Actions"`
}

// Huawei: Oem data for Manager endpoint and SecurityService endpoint
type HttpsCertActionsOemHuawei struct {
	GenerateCSR                LinkTargets  `json:"#HpHttpsCert.GenerateCSR"`
	ImportCertificate          LinkTargets  `json:"#HttpsCert.ImportServerCertificate"`
	X509CertificateInformation X509CertInfo `json:"X509CertificateInformation"`
}

type HttpsCertDataOemHuawei struct {
	CSR     *string                   `json:"CertificateSigningRequest"`
	Id      *string                   `json:"Id"`
	Actions HttpsCertActionsOemHuawei `json:"Actions"`
}

type ManagerDataOemHuaweiLoginRule struct {
	MemberId    *string `json:"MemberId"`
	RuleEnabled bool    `json:"RuleEnabled"`
	StartTime   *string `json:"StartTime"`
	EndTime     *string `json:"EndTime"`
	IP          *string `json:"IP"`
	Mac         *string `json:"Mac"`
}
type SecurityServiceDataOemHuaweiLinks struct {
	HttpsCert OData `json:"HttpsCert"`
}

type SecurityServiceDataOemHuawei struct {
	Id    *string                           `json:"Id"`
	Name  *string                           `json:"Name"`
	Links SecurityServiceDataOemHuaweiLinks `json:"Links"`
}

type _managerDataOemHuawei struct {
	BMCUpTime       *string                         `json:"BMCUpTime"`
	ProductUniqueID *string                         `json:"ProductUniqueID"`
	PlatformType    *string                         `json:"PlatformType"`
	LoginRule       []ManagerDataOemHuaweiLoginRule `json:"LoginRule"`

	SecurityService OData `json:"SecurityService"`
	SnmpService     OData `json:"SnmpService"`
	SmtpService     OData `json:"SmtpService"`
	SyslogService   OData `json:"SyslogService"`
	KvmService      OData `json:"KvmService"`
	NtpService      OData `json:"NtpService"`
}

type ManagerDataOemHuawei struct {
	Huawei _managerDataOemHuawei `json:"Huawei"`
}

// Supermicro: Oem data for Manager.Actions endpoint
type ManagerActionsDataOemSupermicro struct {
	Oem _managerActionsDataOemSupermicro `json:"Oem"`
}

type _managerActionsDataOemSupermicro struct {
	ManagerReset LinkTargets `json:"#Manager.Reset"`
}

// data for CSR subject
type CSRData struct {
	C  string // Country
	S  string // State or province
	L  string // Locality or city
	O  string // Organisation
	OU string // Organisational unit
	CN string // Common name
}

const (
	REDFISH_GENERAL uint = iota
	REDFISH_HP
	REDFISH_HUAWEI
	REDFISH_INSPUR
	REDFISH_LENOVO
	REDFISH_SUPERMICRO
)

type HttpResult struct {
	Url        string
	StatusCode int
	Status     string
	Header     http.Header
	Content    []byte
}

type BaseRedfish interface {
	Initialize() error
	Login() error
	Logout() error
	GetSystems() ([]string, error)
	GetSystemData(string) (*SystemData, error)
	MapSystensById() (map[string]*SystemData, error)
	MapSystemsByUuid() (map[string]*SystemData, error)
	MapSystemsBySerialNumber() (map[string]*SystemData, error)
	GetAccounts() ([]string, error)
	GetAccountData(string) (*AccountData, error)
	MapAccountsByName() (map[string]*AccountData, error)
	MapAccountsById() (map[string]*AccountData, error)
	GetRoles() ([]string, error)
	GetRoleData(string) (*AccountData, error)
	MapRolesByName() (map[string]*RoleData, error)
	MapRolesById() (map[string]*RoleData, error)
	GenCSR(CSRData) error
	FetchCSR() (string, error)
	ImportCertificate(string) error
	ResetSP() error

	httpRequest(string, string, *map[string]string, io.Reader, bool) (HttpResult, error)
	getCSRTarget_HP(*ManagerData) (string, error)
	getCSRTarget_Huawei(*ManagerData) (string, error)
	makeCSRPayload(CSRData) string
	makeCSRPayload_HP(CSRData) string
	makeCSRPayload_Vanilla(CSRData) string
	fetchCSR_HP(*ManagerData) (string, error)
	fetchCSR_Huawei(*ManagerData) (string, error)
	getImportCertTarget_HP(*ManagerData) (string, error)
	getImportCertTarget_Huawei(*ManagerData) (string, error)
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
