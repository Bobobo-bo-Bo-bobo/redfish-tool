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

func (r *Redfish) getImportCertTarget_HP(mgr *ManagerData) (string, error) {
	var certTarget string
	var oemHp ManagerDataOemHp
	var secsvc string
	var oemSSvc SecurityServiceDataOemHp
	var url string
	var transp *http.Transport
	var httpscertloc string
	var httpscert HttpsCertDataOemHp

	// parse Oem section from JSON
	err := json.Unmarshal(mgr.Oem, &oemHp)
	if err != nil {
		return certTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.Id == nil {
		return certTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHp.Hp.Links.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return certTarget, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
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

	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, secsvc)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, secsvc)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return certTarget, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return certTarget, err
	}
	response.Close = true

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return certTarget, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return certTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, httpscertloc)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, httpscertloc)
	}
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return certTarget, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err = client.Do(request)
	if err != nil {
		return certTarget, err
	}
	response.Close = true

	// store unparsed content
	raw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return certTarget, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return certTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return certTarget, err
	}

	if httpscert.Actions.ImportCertificate.Target == nil {
		return certTarget, errors.New(fmt.Sprintf("BUG: .Actions.ImportCertificate.Target is not present or empty in JSON data from %s", url))
	}

	certTarget = *httpscert.Actions.ImportCertificate.Target
	return certTarget, nil
}

func (r *Redfish) getImportCertTarget_Huawei(mgr *ManagerData) (string, error) {
	var certTarget string

	return certTarget, nil
}

func (r *Redfish) ImportCertificate(cert string) error {
	var certtarget string = ""
	var transp *http.Transport
	var url string

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	// get list of Manager endpoint
	mgr_list, err := r.GetManagers()
	if err != nil {
		return err
	}

	// pick the first entry
	mgr0, err := r.GetManagerData(mgr_list[0])
	if err != nil {
		return err
	}

	// get endpoint SecurityService from Managers
	if r.Flavor == REDFISH_HP {
		certtarget, err = r.getImportCertTarget_HP(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == REDFISH_HUAWEI {
	} else if r.Flavor == REDFISH_INSPUR {
		return errors.New("ERROR: Inspur management boards do not support certificate import")
	} else if r.Flavor == REDFISH_SUPERMICRO {
		return errors.New("ERROR: SuperMicro management boards do not support certificate import")
	} else {
		return errors.New("ERROR: Unable to get vendor for management board. If this vendor supports certificate import please file a feature request")
	}

	if certtarget == "" {
		return errors.New("BUG: Target for certificate import is not known")
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
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

	// get URL for Systems endpoint
	client := &http.Client{
		Timeout:   r.Timeout,
		Transport: transp,
	}
	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, certtarget)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, certtarget)
	}

	// escape new lines
	rawcert := strings.Replace(cert, "\n", "\\n", -1)
	cert_payload := fmt.Sprintf("{ \"Certificate\": \"%s\" }", rawcert)

	request, err := http.NewRequest("POST", url, strings.NewReader(cert_payload))
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

	defer request.Body.Close()
	defer response.Body.Close()

	// XXX: do we need to look at the content returned by HTTP POST ?
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("ERROR: HTTP POST to %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	return nil
}
