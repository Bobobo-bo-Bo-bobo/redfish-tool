package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"io/ioutil"
)

func ImportCertificate(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("import-cert", flag.ExitOnError)

	var pem = argParse.String("certificate", "", "Certificate file in PEM format to import")

	argParse.Parse(args)

	if *pem == "" {
		return errors.New("ERROR: Missing mandatory parameter -certificate")
	}

	raw_pem, err := ioutil.ReadFile(*pem)
	if err != nil {
		return err
	}

	// Initialize session
	err = r.Initialise()
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Initialisation failed for %s: %s\n", r.Hostname, err.Error()))
	}

	// Login
	err = r.Login()
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Login to %s failed: %s\n", r.Hostname, err.Error()))
	}

	defer r.Logout()

	// check if vendor support roles
	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HAS_SECURITYSERVICE != redfish.HAS_SECURITYSERVICE {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support certificate import")
		}
	}

	err = r.ImportCertificate(string(raw_pem))
	if err != nil {
		return err
	}

	return nil
}
