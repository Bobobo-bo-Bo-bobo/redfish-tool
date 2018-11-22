package redfish

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

// Initialise Redfish basic data
func (r *Redfish) Initialise(cfg *RedfishConfiguration) error {
    var url string
    var base baseEndpoint

    // get URL for SessionService endpoint
    client := &http.Client{
        Timeout: cfg.Timeout,
    }
    if cfg.Port > 0 {
        url = fmt.Sprintf("http://%s:%d/redfish/v1/", cfg.Hostname, cfg.Port)
    } else {
        url = fmt.Sprintf("http://%s/redfish/v1/", cfg.Hostname)
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
        return errors.New(fmt.Sprintf("BUG: No AccountService endpoint found in base configuration from %s" , url))
    }
    cfg.accountService = *base.AccountService.Id

    if base.Chassis.Id == nil {
        return errors.New(fmt.Sprintf("BUG: No Chassis endpoint found in base configuration from %s" , url))
    }
    cfg.chassis = *base.Chassis.Id

    if base.Managers.Id == nil {
        return errors.New(fmt.Sprintf("BUG: No Managers endpoint found in base configuration from %s" , url))
    }
    cfg.managers = *base.Managers.Id

    if base.SessionService.Id == nil {
        return errors.New(fmt.Sprintf("BUG: No SessionService endpoint found in base configuration from %s", url))
    }
    cfg.sessionService = *base.SessionService.Id

    if base.Systems.Id == nil {
        return errors.New(fmt.Sprintf("BUG: No Systems endpoint found in base configuration from %s" , url))
    }
    cfg.systems = *base.Systems.Id

    return nil
}

// Login to SessionEndpoint and get authentication token for this session
func (r *Redfish) Login(cfg *RedfishConfiguration) error {
    var url string

    if cfg.Username == "" || cfg.Password == "" {
        return errors.New(fmt.Sprintf("ERROR: Both Username and Password must be set"))
    }

    jsonPayload := fmt.Sprintf("{ \"UserName\":\"%s\",\"Password\":\"%s\" }", cfg.Username, cfg.Password)

    // get URL for SessionService endpoint
    client := &http.Client{
        Timeout: cfg.Timeout,
    }
    if cfg.Port > 0 {
        url = fmt.Sprintf("http://%s:%d%s", cfg.Hostname, cfg.Port, cfg.sessionService)
    } else {
        url = fmt.Sprintf("http://%s%s", cfg.Hostname, cfg.sessionService)
    }

    request, err := http.NewRequest("POST", url, strings.NewReader(jsonPayload))
    if err != nil {
        return err
    }

    request.Header.Add("OData-Version", "4.0")
    request.Header.Add("Accept", "application/json")
    request.Header.Add("Content-Type", "application/json")

    response, err := client.Do(request)
    if err != nil {
        return err
    }
    
    defer response.Body.Close()

    if response.StatusCode != 200 {
        return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
    }

    token := response.Header.Get("x-auth-token")
    if token == "" {
        return errors.New(fmt.Sprintf("BUG: HTTP POST to SessionService endpoint %s returns OK but no X-Auth-Token in reply", url))
    }

    cfg.AuthToken = &token
    return nil
}

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
        url = fmt.Sprintf("http://%s:%d%s", cfg.Hostname, cfg.Port, cfg.sessionService)
    } else {
        url = fmt.Sprintf("http://%s%s", cfg.Hostname, cfg.sessionService)
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
        return errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
    }

    cfg.AuthToken = nil

    return nil
}

