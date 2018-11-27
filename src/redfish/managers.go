package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//get array of managers and their endpoints
func (r *Redfish) GetManagers() ([]string, error) {
	var url string
	var mgrs OData
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
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, r.Managers)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, r.Managers)
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

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return result, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &mgrs)
	if err != nil {
		return result, err
	}

	if len(mgrs.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Managers"))
	}

	for _, m := range mgrs.Members {
		result = append(result, *m.Id)
	}
	return result, nil
}

// get manager data for a particular account
func (r *Redfish) GetManagerData(managerEndpoint string) (*ManagerData, error) {
	var result ManagerData
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
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, managerEndpoint)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, managerEndpoint)
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
	result.SelfEndpoint = &managerEndpoint
	return &result, nil
}

// map ID -> manager data
func (r *Redfish) MapManagersById() (map[string]*ManagerData, error) {
	var result = make(map[string]*ManagerData)

	ml, err := r.GetManagers()
	if err != nil {
		return result, err
	}

	for _, mgr := range ml {
		m, err := r.GetManagerData(mgr)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if m.Id == nil {
			return result, errors.New(fmt.Sprintf("BUG: No Id found or Id is null in JSON data from %s", mgr))
		}
		result[*m.Id] = m
	}

	return result, nil
}

// map UUID -> manager data
func (r *Redfish) MapManagersByUuid() (map[string]*ManagerData, error) {
	var result = make(map[string]*ManagerData)

	ml, err := r.GetManagers()
	if err != nil {
		return result, err
	}

	for _, mgr := range ml {
		m, err := r.GetManagerData(mgr)
		if err != nil {
			return result, err
		}

		// should NEVER happen
		if m.UUID == nil {
			return result, errors.New(fmt.Sprintf("BUG: No UUID found or UUID is null in JSON data from %s", mgr))
		}
		result[*m.UUID] = m
	}

	return result, nil
}
