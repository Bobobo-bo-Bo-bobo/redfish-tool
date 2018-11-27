package main

import (
	"errors"
	"fmt"

	"redfish"
)

func GetAllManagers(r redfish.Redfish) error {
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
	// get all manager endpoints
	mmap, err := r.MapManagersById()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	// loop over all endpoints
	for mname, mgr := range mmap {

		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + mname)
		if mgr.Id != nil {
			fmt.Println("  Id: " + *mgr.Id)
		}
		if mgr.Name != nil {
			fmt.Println("  Name:", *mgr.Name)
		}

		if mgr.ManagerType != nil {
			fmt.Println("  ManagerType:", *mgr.ManagerType)
		}

		if mgr.UUID != nil {
			fmt.Println("  UUID:", *mgr.UUID)
		}

		if mgr.FirmwareVersion != nil {
			fmt.Println("  FirmwareVersion:", *mgr.FirmwareVersion)
		}

		fmt.Println("  Status: ")
		if mgr.Status.State != nil {
			fmt.Println("   State: " + *mgr.Status.State)
		}
		if mgr.Status.Health != nil {
			fmt.Println("   Health: " + *mgr.Status.Health)
		}
		if mgr.SelfEndpoint != nil {
			fmt.Println("  Endpoint: " + *mgr.SelfEndpoint)
		}

	}

	return nil
}
