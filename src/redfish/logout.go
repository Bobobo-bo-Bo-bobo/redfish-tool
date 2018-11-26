package redfish

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

// Logout from SessionEndpoint and delete authentication token for this session
func (r *Redfish) Logout() error {
	var url string
	var transp *http.Transport

	if r.AuthToken == nil {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}
	if *r.AuthToken == "" {
		// do nothing for Logout when we don't even have an authentication token
		return nil
	}

	if r.SessionLocation == nil || *r.SessionLocation == "" {
		return errors.New(fmt.Sprintf("BUG: X-Auth-Token set (value: %s) but no SessionLocation for this session found\n", *r.AuthToken))
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
	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}

	url = *r.SessionLocation

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	response.Close = true

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP DELETE for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	r.AuthToken = nil
	r.SessionLocation = nil

	return nil
}
