package main

import (
	"errors"
	"fmt"

	redfish "git.ypbind.de/repository/go-redfish.git"
)

func GetAllSystems(r redfish.Redfish) error {
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

	// get all systems
	smap, err := r.MapSystemsById()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	for sname, sdata := range smap {
		fmt.Println(" " + sname)

		if sdata.Id != nil {
			fmt.Println("  Id:", *sdata.Id)
		}

		if sdata.UUID != nil {
			fmt.Println("  UUID:", *sdata.UUID)
		}

		if sdata.Name != nil {
			fmt.Println("  Name:", *sdata.Name)
		}

		if sdata.SerialNumber != nil {
			fmt.Println("  SerialNumber:", *sdata.SerialNumber)
		}

		if sdata.Manufacturer != nil {
			fmt.Println("  Manufacturer:", *sdata.Manufacturer)
		}

		if sdata.Model != nil {
			fmt.Println("  Model:", *sdata.Model)
		}

		fmt.Println("  Status:")
		if sdata.Status.State != nil {
			fmt.Println("   State: " + *sdata.Status.State)
		}
		if sdata.Status.Health != nil {
			fmt.Println("   Health: " + *sdata.Status.Health)
		}
		if sdata.Status.HealthRollUp != nil {
			fmt.Println("   HealthRollUp: " + *sdata.Status.HealthRollUp)
		}

		if sdata.PowerState != nil {
			fmt.Println("  PowerState:", *sdata.PowerState)
		}

		if sdata.BIOSVersion != nil {
			fmt.Println("  BIOSVersion:", *sdata.BIOSVersion)
		}

		if sdata.SelfEndpoint != nil {
			fmt.Println("  SelfEndpoint:", *sdata.SelfEndpoint)
		}

	}
	return nil
}
