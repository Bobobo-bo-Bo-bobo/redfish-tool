package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//get array of roles and their endpoints
func (r *Redfish) GetRoles(cfg *RedfishConfiguration) ([]string, error) {
	var url string
	var accsvc AccountService
	var roles OData
	var transp *http.Transport
	var result = make([]string, 0)

	if cfg.AuthToken == nil || *cfg.AuthToken == "" {
		return result, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	if cfg.InsecureSSL {
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
		Timeout:   cfg.Timeout,
		Transport: transp,
	}
	if cfg.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", cfg.Hostname, cfg.Port, cfg.AccountService)
	} else {
		url = fmt.Sprintf("https://%s%s", cfg.Hostname, cfg.AccountService)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *cfg.AuthToken)

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

	err = json.Unmarshal(raw, &accsvc)
	if err != nil {
		return result, err
	}

	// Some managementboards (e.g. HPE iLO) don't use roles but an internal ("Oem") privilege map instead
	if accsvc.RolesEndpoint == nil {
		return result, nil
	}

	if cfg.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", cfg.Hostname, cfg.Port, *accsvc.RolesEndpoint.Id)
	} else {
		url = fmt.Sprintf("https://%s%s", cfg.Hostname, *accsvc.RolesEndpoint.Id)
	}
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *cfg.AuthToken)

	request.Close = true

	response, err = client.Do(request)
	if err != nil {
		return result, err
	}
	response.Close = true

	// store unparsed content
	raw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return result, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return result, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &roles)
	if err != nil {
		return result, err
	}

	if len(roles.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Roles"))
	}

	for _, r := range roles.Members {
		result = append(result, *r.Id)
	}
	return result, nil
}

// get role data for a particular role
func (r *Redfish) GetRoleData(cfg *RedfishConfiguration, roleEndpoint string) (*RoleData, error) {
	var result RoleData
	var url string
	var transp *http.Transport

	if cfg.AuthToken == nil || *cfg.AuthToken == "" {
		return nil, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
	}

	if cfg.InsecureSSL {
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
		Timeout:   cfg.Timeout,
		Transport: transp,
	}
	if cfg.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", cfg.Hostname, cfg.Port, roleEndpoint)
	} else {
		url = fmt.Sprintf("https://%s%s", cfg.Hostname, roleEndpoint)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *cfg.AuthToken)

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

	return &result, nil
}
