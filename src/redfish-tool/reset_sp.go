package main

import (
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func resetSP(r redfish.Redfish) error {
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

	fmt.Println(r.Hostname)

	err = r.ResetSP()
	if err != nil {
		return err
	}

	return nil
}
