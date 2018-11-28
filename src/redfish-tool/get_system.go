package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"redfish"
)

func GetSystem(r redfish.Redfish, args []string) error {
	var sys *redfish.SystemData
	var found bool
	var smap map[string]*redfish.SystemData

	argParse := flag.NewFlagSet("get-system", flag.ExitOnError)

	var uuid = argParse.String("uuid", "", "Get detailed information for system identified by UUID")
	var id = argParse.String("id", "", "Get detailed information for system identified by ID")

	argParse.Parse(args)

	if *uuid != "" && *id != "" {
		return errors.New("ERROR: Options -uuid and -id are mutually exclusive")
	}

	if *uuid == "" && *id == "" {
		return errors.New("ERROR: Required options -uuid or -id not found")
	}

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
	if *id != "" {
		smap, err = r.MapSystemsById()
	} else {
		smap, err = r.MapSystemsByUuid()
	}

	if err != nil {
		return err
	}

	if *id != "" {
		sys, found = smap[*id]
	} else {
		sys, found = smap[*uuid]
	}

	fmt.Println(r.Hostname)
	if found {
		if *id != "" {
			fmt.Println(" " + *sys.Id)
		} else {
			fmt.Println(" " + *sys.UUID)
		}

		if sys.Id != nil {
			fmt.Println("  Id:", *sys.Id)
		}

		if sys.UUID != nil {
			fmt.Println("  UUID:", *sys.UUID)
		}

		if sys.Name != nil {
			fmt.Println("  Name:", *sys.Name)
		}

		if sys.SerialNumber != nil {
			fmt.Println("  SerialNumber:", *sys.SerialNumber)
		}

		if sys.Manufacturer != nil {
			fmt.Println("  Manufacturer:", *sys.Manufacturer)
		}

		if sys.Model != nil {
			fmt.Println("  Model:", *sys.Model)
		}

		fmt.Println("  Status:")
		if sys.Status.State != nil {
			fmt.Println("   State: " + *sys.Status.State)
		}
		if sys.Status.Health != nil {
			fmt.Println("   Health: " + *sys.Status.Health)
		}
		if sys.Status.HealthRollUp != nil {
			fmt.Println("   HealthRollUp: " + *sys.Status.HealthRollUp)
		}

		if sys.PowerState != nil {
			fmt.Println("  PowerState:", *sys.PowerState)
		}

		if sys.BIOSVersion != nil {
			fmt.Println("  BIOSVersion:", *sys.BIOSVersion)
		}

		if sys.SelfEndpoint != nil {
			fmt.Println("  SelfEndpoint:", *sys.SelfEndpoint)
		}
	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *uuid, r.Hostname)
		}

	}
	return nil
}
