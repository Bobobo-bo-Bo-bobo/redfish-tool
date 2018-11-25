package redfish

type oData struct {
	Id           *string `json:"@odata.id"`
	Type         *string `json:"@odata.type"`
	Context      *string `json:"@odata.context"`
	Members      []oData `json:"Members"`
	MembersCount int     `json:"Members@odata.count"`
}

type baseEndpoint struct {
	AccountService oData `json:"AccountService"`
	Chassis        oData `json:"Chassis"`
	Managers       oData `json:"Managers"`
	SessionService oData `json:"SessionService"`
	Systems        oData `json:"Systems"`
}

type sessionServiceEndpoint struct {
	Enabled        *bool  `json:"ServiceEnabled"`
	SessionTimeout int    `json:"SessionTimeout"`
	Sessions       *oData `json:"Sessions"`
}
