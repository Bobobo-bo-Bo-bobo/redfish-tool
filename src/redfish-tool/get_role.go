package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"os"
)

func GetRole(r redfish.Redfish, args []string) error {
	var rle *redfish.RoleData
	var found bool
	var rmap map[string]*redfish.RoleData
	argParse := flag.NewFlagSet("get-role", flag.ExitOnError)

	var id = argParse.String("id", "", "Get detailed information for role identified by ID")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *id == "" {
		return errors.New("ERROR: Required option -id not found")
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

	// check if vendor support roles
	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HAS_ACCOUNT_ROLES != redfish.HAS_ACCOUNT_ROLES {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support roles")
		}
	}

	// get all roles
	rmap, err = r.MapRolesById()

	if err != nil {
		return err
	}

	rle, found = rmap[*id]

	if found {
		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + *id)
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

	} else {
		fmt.Fprintf(os.Stderr, "Role %s not found on %s\n", *id, r.Hostname)
	}

	return nil
}
