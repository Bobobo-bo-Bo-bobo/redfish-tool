package redfish

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

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

