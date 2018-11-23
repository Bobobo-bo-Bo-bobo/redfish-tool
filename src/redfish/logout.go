package redfish

import (
	//    "encoding/json"
	"errors"
	"fmt"
	//    "io/ioutil"
	"net/http"
	//    "strings"
)

// Logout from SessionEndpoint and delete authentication token for this session
func (r *Redfish) Logout(cfg *RedfishConfiguration) error {
	var url string

	if cfg.AuthToken == nil {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	if cfg.Port > 0 {
		url = fmt.Sprintf("http://%s:%d%s", cfg.Hostname, cfg.Port, cfg.sessions)
	} else {
		url = fmt.Sprintf("http://%s%s", cfg.Hostname, cfg.sessions)
	}

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Token", *cfg.AuthToken)

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("ERROR: HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	cfg.AuthToken = nil

	return nil
}
