package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"os"
)

func GetLicense(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("get-license", flag.ExitOnError)
	var id = argParse.String("id", "", "Management board identified by ID")
	var uuid = argParse.String("uuid", "", "Management board identified by UUID")
	var mmap map[string]*redfish.ManagerData
	var mgr *redfish.ManagerData

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

	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HAS_LICENSE != redfish.HAS_LICENSE {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support license operations")
		}
	}

	if *id != "" {
		mmap, err = r.MapManagersById()
	} else if *uuid != "" {
		mmap, err = r.MapManagersByUuid()
	}

	if err != nil {
		return err
	}

	if *id != "" {
		mgr, found = mmap[*id]
	} else if *uuid != "" {
		mgr, found = mmap[*uuid]
	}

	fmt.Println(r.Hostname)
	if found {
		l, err := r.GetLicense(mgr)
		if err != nil {
			return err
		}

		if l.Name != "" {
			fmt.Println(" Name: " + l.Name)
		} else {
			fmt.Println(" Name: -")
		}

		if l.Type != "" {
			fmt.Println(" Type: " + l.Type)
		} else {
			fmt.Println(" Type: -")
		}

		if l.Expiration != "" {
			fmt.Println(" Expiration: " + l.Expiration)
		} else {
			fmt.Println(" Expiration: -")
		}

		if l.License != "" {
			fmt.Println(" License: " + l.License)
		} else {
			fmt.Println(" License: -")
		}

	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "Manager with ID %s not found on %s\n", *id, r.Hostname)
		} else if *uuid != "" {
			fmt.Fprintf(os.Stderr, "Manager with UUID %s not found on %s\n", *uuid, r.Hostname)
		}
	}

	return nil
}
