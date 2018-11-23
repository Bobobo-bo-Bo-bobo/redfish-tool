package redfish

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Initialise Redfish basic data
func (r *Redfish) Initialise(cfg *RedfishConfiguration) error {
	var url string
	var base baseEndpoint
	var transp *http.Transport

	if cfg.InsecureSSL {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		transp = &http.Transport{
			TLSClientConfig: &tls.Config{},
		}
	}

	// get URL for SessionService endpoint
	client := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transp,
	}
	if cfg.Port > 0 {
		url = fmt.Sprintf("https://%s:%d/redfish/v1/", cfg.Hostname, cfg.Port)
	} else {
		url = fmt.Sprintf("https://%s/redfish/v1/", cfg.Hostname)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	cfg.rawBaseContent = string(raw)

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &base)
	if err != nil {
		return err
	}

	// extract required endpoints
	if base.AccountService.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No AccountService endpoint found in base configuration from %s", url))
	}
	cfg.accountService = *base.AccountService.Id

	if base.Chassis.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Chassis endpoint found in base configuration from %s", url))
	}
	cfg.chassis = *base.Chassis.Id

	if base.Managers.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Managers endpoint found in base configuration from %s", url))
	}
	cfg.managers = *base.Managers.Id

	if base.SessionService.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No SessionService endpoint found in base configuration from %s", url))
	}
	cfg.sessionService = *base.SessionService.Id

	if base.Systems.Id == nil {
		return errors.New(fmt.Sprintf("BUG: No Systems endpoint found in base configuration from %s", url))
	}
	cfg.systems = *base.Systems.Id

	return nil
}
