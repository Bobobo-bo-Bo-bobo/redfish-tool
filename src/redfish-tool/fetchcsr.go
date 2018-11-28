package main

import (
	"errors"
	"fmt"
	"redfish"
)

func FetchCSR(r redfish.Redfish) error {
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

	csr, err := r.FetchCSR()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	fmt.Println(csr)

	return nil
}