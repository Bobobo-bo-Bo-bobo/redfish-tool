package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//get array of systems and their endpoints
func (r *Redfish) GetSystems() ([]string, error) {
	var url string
	var systems OData
	var transp *http.Transport
	var result = make([]string, 0)

	if r.AuthToken == nil || *r.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	if r.InsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// get URL for Systems endpoint
	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}
	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, r.Systems)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, r.Systems)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return result, err
	}
	response.Close = true

	defer response.Body.Close()

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &systems)
	if err != nil {
		return result, err
	}

	if len(systems.Members) == 0 {
		return result, errors.New("BUG: Array of system endpoints is empty")
	}

	for _, endpoint := range systems.Members {
		result = append(result, *endpoint.Id)
	}
	return result, nil
}

// get system data for a particular system
func (r *Redfish) GetSystemData(systemEndpoint string) (*SystemData, error) {
	var result SystemData
	var url string
	var transp *http.Transport

	if r.AuthToken == nil || *r.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	if r.InsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// get URL for Systems endpoint
	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}
	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, systemEndpoint)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, systemEndpoint)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	response.Close = true

	defer response.Body.Close()

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	result.SelfEndpoint = &systemEndpoint
	return &result, nil
}

// Map systems by ID
func (r *Redfish) MapSystemsById() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found for System at %s", sys))
		}

		result[*s.Id] = s
	}

	return result, nil
}

// Map systems by UUID
func (r *Redfish) MapSystemsByUuid() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.UUID == nil {
			return result, errors.New(fmt.Sprintf("BUG: No UUID found for System at %s", sys))
		}

		result[*s.UUID] = s
	}

	return result, nil
}

// Map systems by serial number
func (r *Redfish) MapSystemsBySerialNumber() (map[string]*SystemData, error) {
	var result = make(map[string]*SystemData)

	sysl, err := r.GetSystems()
	if err != nil {
		return result, nil
	}

	for _, sys := range sysl {
		s, err := r.GetSystemData(sys)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if s.SerialNumber == nil {
			return result, errors.New(fmt.Sprintf("BUG: No SerialNumber found for System at %s", sys))
		}

		result[*s.SerialNumber] = s
	}

	return result, nil
}

// get vendor specific "flavor"
func (r *Redfish) GetVendorFlavor() error {
	// get vendor "flavor" for vendor specific implementation details
	_sys, err := r.GetSystems()
	if err != nil {
		return err
	}
	// assuming every system has the same vendor, pick the first one to determine vendor flavor
	_sys0, err := r.GetSystemData(_sys[0])
	if _sys0.Manufacturer != nil {
		_manufacturer := strings.TrimSpace(strings.ToLower(*_sys0.Manufacturer))
		if _manufacturer == "hp" || _manufacturer == "hpe" {
			r.Flavor = REDFISH_HP
		} else if _manufacturer == "huawei" {
			r.Flavor = REDFISH_HUAWEI
		} else if _manufacturer == "inspur" {
			r.Flavor = REDFISH_INSPUR
		} else if _manufacturer == "supermicro" {
			r.Flavor = REDFISH_SUPERMICRO
		} else {
			r.Flavor = REDFISH_GENERAL
		}
	}

	return nil
}
