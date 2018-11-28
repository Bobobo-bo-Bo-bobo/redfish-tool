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

func (r *Redfish) getCSRTarget_HP(mgr *ManagerData) (string, error) {
	var csrTarget string
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
		return csrTarget, err
	}

	// get SecurityService endpoint from .Oem.Hp.links.SecurityService
	if oemHp.Hp.Links.SecurityService.Id == nil {
		return csrTarget, errors.New("BUG: .Hp.Links.SecurityService.Id not found or null")
	} else {
		secsvc = *oemHp.Hp.Links.SecurityService.Id
	}

	if r.AuthToken == nil || *r.AuthToken == "" {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: No authentication token found, is the session setup correctly?"))
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
		return csrTarget, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		return csrTarget, err
	}
	response.Close = true

	// store unparsed content
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return csrTarget, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &oemSSvc)
	if err != nil {
		return csrTarget, err
	}

	if oemSSvc.Links.HttpsCert.Id == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .links.HttpsCert.Id not present or is null in data from %s", url))
	}

	httpscertloc = *oemSSvc.Links.HttpsCert.Id

	if r.Port > 0 {
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, httpscertloc)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, httpscertloc)
	}
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return csrTarget, err
	}

	request.Header.Add("OData-Version", "4.0")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("X-Auth-Token", *r.AuthToken)

	request.Close = true

	response, err = client.Do(request)
	if err != nil {
		return csrTarget, err
	}
	response.Close = true

	// store unparsed content
	raw, err = ioutil.ReadAll(response.Body)
	if err != nil {
		response.Body.Close()
		return csrTarget, err
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return csrTarget, errors.New(fmt.Sprintf("ERROR: HTTP GET for %s returned \"%s\" instead of \"200 OK\"", url, response.Status))
	}

	err = json.Unmarshal(raw, &httpscert)
	if err != nil {
		return csrTarget, err
	}

	if httpscert.Actions.GenerateCSR.Target == nil {
		return csrTarget, errors.New(fmt.Sprintf("BUG: .Actions.GenerateCSR.Target is not present or empty in JSON data from %s", url))
	}

	csrTarget = *httpscert.Actions.GenerateCSR.Target
	return csrTarget, nil
}

func (r *Redfish) getCSRTarget_Huawei(mgr *ManagerData) (string, error) {
	var csrTarget string

	return csrTarget, nil
}

func (r *Redfish) GenCSR(csr CSRData) error {
	var csrstr string = ""
	var gencsrtarget string = ""
	var transp *http.Transport
	var url string

	// set vendor flavor
	err := r.GetVendorFlavor()
	if err != nil {
		return err
	}

	if csr.C != "" {
		csrstr += fmt.Sprintf("\"Country\": \"%s\", ", csr.C)
	}

	if csr.S != "" {
		csrstr += fmt.Sprintf("\"State\": \"%s\", ", csr.S)
	}

	if csr.L != "" {
		csrstr += fmt.Sprintf("\"City\": \"%s\", ", csr.L)
	}

	if csr.O != "" {
		csrstr += fmt.Sprintf("\"OrgName\": \"%s\", ", csr.O)
	}

	if csr.OU != "" {
		csrstr += fmt.Sprintf("\"OrgUnit\": \"%s\", ", csr.OU)
	}

	if csr.CN != "" {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", csr.CN)
	} else {
		csrstr += fmt.Sprintf("\"CommonName\": \"%s\" ", r.Hostname)
	}

	csrstr = "{ " + csrstr + " } "

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
		gencsrtarget, err = r.getCSRTarget_HP(mgr0)
		if err != nil {
			return err
		}
	} else if r.Flavor == REDFISH_HUAWEI {
	} else if r.Flavor == REDFISH_INSPUR {
		return errors.New("ERROR: Inspur management boards do not support CSR generation")
	} else if r.Flavor == REDFISH_SUPERMICRO {
		return errors.New("ERROR: SuperMicro management boards do not support CSR generation")
	} else {
		return errors.New("ERROR: Unable to get vendor for management board. If this vendor supports CSR generation please file a feature request")
	}

	if gencsrtarget == "" {
		return errors.New("BUG: CSR generation target is not known")
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
		url = fmt.Sprintf("https://%s:%d%s", r.Hostname, r.Port, gencsrtarget)
	} else {
		url = fmt.Sprintf("https://%s%s", r.Hostname, gencsrtarget)
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(csrstr))
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
