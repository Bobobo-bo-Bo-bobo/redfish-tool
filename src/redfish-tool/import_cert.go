package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"io/ioutil"
	"os"
)

func importCertificate(r redfish.Redfish, args []string) error {
	var rawPem []byte
	var err error

	argParse := flag.NewFlagSet("import-cert", flag.ExitOnError)

	var pem = argParse.String("certificate", "", "Certificate file in PEM format to import")

	argParse.Parse(args)

	if *pem == "" {
		return errors.New("ERROR: Missing mandatory parameter -certificate")
	}

	if *pem == "-" {
		rawPem, err = ioutil.ReadAll(os.Stdin)
	} else {
		rawPem, err = ioutil.ReadFile(*pem)
	}
	if err != nil {
		return err
	}

	// Initialize session
	err = r.Initialise()
	if err != nil {
		return fmt.Errorf("ERROR: Initialisation failed for %s: %s", r.Hostname, err.Error())
	}

	// Login
	err = r.Login()
	if err != nil {
		return fmt.Errorf("ERROR: Login to %s failed: %s", r.Hostname, err.Error())
	}

	defer r.Logout()

	// check if vendor support roles
	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HasSecurityService != redfish.HasSecurityService {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support certificate import")
		}
	}

	err = r.ImportCertificate(string(rawPem))
	if err != nil {
		return err
	}

	return nil
}
