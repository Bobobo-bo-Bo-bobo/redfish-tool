package redfish

import "time"

type RedfishConfiguration struct {
	Hostname       string
	Port           int
	Username       string
	Password       string
	AuthToken      *string
	Timeout        time.Duration
	InsecureSSL    bool
	rawBaseContent string

	// endpoints
	accountService string
	chassis        string
	managers       string
	sessionService string
	sessions       string
	systems        string
}

type Result struct {
	RawContent *string
	RawHeaders *string
}

type BaseRedfish interface {
	Initialize(*RedfishConfiguration) error
	Login(*RedfishConfiguration) error
	Logout(*RedfishConfiguration) error
}

type Redfish struct {
}
