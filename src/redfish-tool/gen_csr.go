package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func GenCSR(r redfish.Redfish, args []string) error {
	var csrdata redfish.CSRData

	argParse := flag.NewFlagSet("gen-csr", flag.ExitOnError)

	var c = argParse.String("country", "", "CSR - country")
	var s = argParse.String("state", "", "CSR - state or province")
	var l = argParse.String("locality", "", "CSR - locality or city")
	var o = argParse.String("organisation", "", "CSR - organisation")
	var ou = argParse.String("organisational-unit", "", "CSR - organisational unit")
	var cn = argParse.String("common-name", "", "CSR - common name")

	argParse.Parse(args)

	// Initialize session
	err := r.Initialise()
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
			return errors.New("Vendor does not support CSR generation")
		}
	}

	csrdata = redfish.CSRData{
		C:  *c,
		S:  *s,
		L:  *l,
		O:  *o,
		OU: *ou,
		CN: *cn,
	}

	err = r.GenCSR(csrdata)
	if err != nil {
		return err
	}

	return nil
}
