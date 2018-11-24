package redfish

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

// Logout from SessionEndpoint and delete authentication token for this session
func (r *Redfish) Logout(cfg *RedfishConfiguration) error {
	var url string
	var transp *http.Transport

	if cfg.AuthToken == nil {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}
	if *cfg.AuthToken == "" {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}

	if cfg.SessionLocation == nil || *cfg.SessionLocation == "" {
		return errors.New(fmt.Sprintf("BUG: X-Auth-Token set (value: %s) but no SessionLocation for this session found\n", *cfg.AuthToken))
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
	client := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transp,
	}

    url = *cfg.SessionLocation

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
