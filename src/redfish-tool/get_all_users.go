package main

import (
	"errors"
	"fmt"

	"redfish"
)

func GetAllUsers(r redfish.Redfish) error {
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
	// get all account endpoints
	amap, err := r.MapAccountNames()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	// loop over all endpoints
	for aname, acc := range amap {

		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + aname)
		if acc.Id != nil && *acc.Id != "" {
			fmt.Println("  Id: " + *acc.Id)
		}

		if acc.Name != nil && *acc.Name != "" {
			fmt.Println("  Name: " + *acc.Name)
		}

		if acc.UserName != nil && *acc.UserName != "" {
			fmt.Println("  UserName: " + *acc.UserName)
		}

		if acc.Password != nil && *acc.Password != "" {
			fmt.Println("  Password: " + *acc.Password)
		}

		if acc.RoleId != nil && *acc.RoleId != "" {
			fmt.Println("  RoleId: " + *acc.RoleId)
		}

		if acc.Enabled != nil {
			if *acc.Enabled {
				fmt.Println("  Enabled: true")
			} else {
				fmt.Println("  Enabled: false")
			}
		}

		if acc.Locked != nil {
			if *acc.Locked {
				fmt.Println("  Locked: true")
			} else {
				fmt.Println("  Locked: false")
			}
		}

        if acc.SelfEndpoint != nil {
            fmt.Println("  Endpoint: " + *acc.SelfEndpoint)
        }
	}

	return nil
}
