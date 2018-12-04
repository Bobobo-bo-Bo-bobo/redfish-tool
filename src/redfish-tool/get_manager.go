package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"os"
)

func GetManager(r redfish.Redfish, args []string) error {
	var mgr *redfish.ManagerData
	var found bool
	var mmap map[string]*redfish.ManagerData
	argParse := flag.NewFlagSet("get-manager", flag.ExitOnError)

	var uuid = argParse.String("uuid", "", "Get detailed information for user identified by UUID")
	var id = argParse.String("id", "", "Get detailed information for user identified by ID")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

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

	// get all account endpoints
	if *id != "" {
		mmap, err = r.MapManagersById()
	} else {
		mmap, err = r.MapManagersByUuid()
	}

	if err != nil {
		return err
	}

	if *id != "" {
		mgr, found = mmap[*id]
	} else {
		mgr, found = mmap[*uuid]
	}

	if found {
		// XXX: Allow for different output formats like JSON, YAML, ... ?
		if *id != "" {
			fmt.Println(" " + *mgr.Id)
		} else {
			fmt.Println(" " + *mgr.UUID)
		}

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

	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *uuid, r.Hostname)
		}
	}

	return nil
}
