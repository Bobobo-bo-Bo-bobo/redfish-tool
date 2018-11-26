package main

import (
	"errors"
	"fmt"

	"redfish"
)

func GetAllRoles(r redfish.Redfish) error {
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

	// get all role endpoints - Note: role names are _NOT_ unique but IDs are!
	rmap, err := r.MapRolesById()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	// loop over all endpoints
	for rid, rle := range rmap {

		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + rid)
		if rle.Id != nil && *rle.Id != "" {
			fmt.Println("  Id: " + *rle.Id)
		}

		if rle.Name != nil && *rle.Name != "" {
			fmt.Println("  Name: " + *rle.Name)
		}

		if rle.IsPredefined != nil {
			if *rle.IsPredefined {
				fmt.Println("  IsPredefined: true")
			} else {
				fmt.Println("  IsPredefined: false")
			}
		}

		if len(rle.AssignedPrivileges) != 0 {
			fmt.Println("  Assigned privieleges")
			for _, p := range rle.AssignedPrivileges {
				fmt.Println("   " + p)
			}
		}

		if rle.SelfEndpoint != nil {
			fmt.Println("  Endpoint: " + *rle.SelfEndpoint)
		}
	}

	return nil
}
