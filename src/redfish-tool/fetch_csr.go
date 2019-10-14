package main

import (
	"errors"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func fetchCSR(r redfish.Redfish) error {
	// Initialize session
	err := r.Initialise()
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
			return errors.New("Vendor does not support CSR generation")
		}
	}

	csr, err := r.FetchCSR()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	fmt.Println(csr)

	return nil
}
