package redfish

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
}

type AccountService struct {
	AccountsEndpoint *OData `json:"Accounts"`
}

type AccountData struct {
	Id       *string `json:"Id"`
	Name     *string `json:"Name"`
	UserName *string `json:"UserName"`
	Password *string `json:"Password"`
	RoleId   *string `json:"RoleId"`
	Enabled  *bool   `json:"Enabled"`
	Locked   *bool   `json:"Locked"`
}
