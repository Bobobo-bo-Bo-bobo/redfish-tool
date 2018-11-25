package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//get array of systems and their endpoints
func (r *Redfish) GetSystems(cfg *RedfishConfiguration) ([]string, error) {
	var url string
	var systems oData
	var transp *http.Transport
	var result = make([]string, 0)

    if cfg.AuthToken == nil || cfg.AuthToken == "" {
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
		url = fmt.Sprintf("https://%s:%d%s", cfg.Hostname, cfg.Port, cfg.Systems)
	} else {
		url = fmt.Sprintf("https://%s%s", cfg.Hostname, cfg.Systems)
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
