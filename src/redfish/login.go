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

// Login to SessionEndpoint and get authentication token for this session
func (r *Redfish) Login() error {
	var url string
	var sessions sessionServiceEndpoint
	var transp *http.Transport

	if r.Username == "" || r.Password == "" {
		return errors.New(fmt.Sprintf("ERROR: Both Username and Password must be set"))
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

	// get URL for SessionService endpoint
	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}

	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, r.SessionService)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, r.SessionService)
	}

	// get Sessions endpoint, which requires HTTP Basic auth
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(r.Username, r.Password)
	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	response.Close = true

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return err
	}

	response.Body.Close()

	err = json.Unmarshal(raw, &sessions)
	if err != nil {
		return err
	}

	// check if management boards reports "ServiceEnabled" and if it does, check if is true
	if sessions.Enabled != nil {
		if !*sessions.Enabled {
			response.Body.Close()
			return errors.New(fmt.Sprintf("ERROR: Session information from %s reports session service as disabled\n", url))
		}
	}

	if sessions.Sessions == nil {
		response.Body.Close()
		return errors.New(fmt.Sprintf("BUG: No Sessions endpoint reported from %s\n", url))
	}

	if sessions.Sessions.Id == nil {
		response.Body.Close()
		return errors.New(fmt.Sprintf("BUG: Malformed Sessions endpoint reported from %s: no @odata.id field found\n", url))
	}

	r.Sessions = *sessions.Sessions.Id

	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, *sessions.Sessions.Id)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, *sessions.Sessions.Id)
	}

	jsonPayload := fmt.Sprintf("{ \"UserName\":\"%s\",\"Password\":\"%s\" }", r.Username, r.Password)
	request, err = http.NewRequest("POST", url, strings.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Close = true

	response, err = client.Do(request)
	if err != nil {
		return err
	}

	response.Close = true
	response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST for %s returned \"%s\" instead of \"200 OK\" or \"201 Created\"", url, response.Status))
	}

	token := response.Header.Get("x-auth-token")
	if token == "" {
		return errors.New(fmt.Sprintf("BUG: HTTP POST to SessionService endpoint %s returns OK but no X-Auth-Token in reply", url))
	}
	r.AuthToken = &token

	session := response.Header.Get("location")
	if session == "" {
		return errors.New(fmt.Sprintf("BUG: HTTP POST to SessionService endpoint %s returns OK but has no Location in reply", url))
	}

	// check if is a full URL
	if session[0] == '/' {
		if r.Port > 0 {
			session = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, session)
		} else {
			session = fmt.Sprintf("https://%s%s", r.Hostname, session)
		}
	}
	r.SessionLocation = &session

	return nil
}
