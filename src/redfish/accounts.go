package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//get array of accounts and their endpoints
func (r *Redfish) GetAccounts(cfg *RedfishConfiguration) ([]string, error) {
	var url string
	var accsvc AccountService
	var accs OData
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

	if accsvc.AccountsEndpoint == nil {
		return result, errors.New("BUG: No Accounts endpoint found")
	}

	if cfg.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", cfg.Hostname, cfg.Port, *accsvc.AccountsEndpoint.Id)
	} else {
		url = fmt.Sprintf("https://%s%s", cfg.Hostname, *accsvc.AccountsEndpoint.Id)
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

	err = json.Unmarshal(raw, &accs)
	if err != nil {
		return result, err
	}

	if len(accs.Members) == 0 {
		return result, errors.New(fmt.Sprintf("BUG: Missing or empty Members attribute in Accounts"))
	}

	for _, a := range accs.Members {
		result = append(result, *a.Id)
	}
	return result, nil
}
