package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"redfish"
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

	err = r.ImportCertificate(string(raw_pem))
	if err != nil {
		return err
	}

	return nil
}
