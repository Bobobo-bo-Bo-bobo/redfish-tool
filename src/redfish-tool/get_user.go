package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"os"
)

func GetUser(r redfish.Redfish, args []string) error {
	var acc *redfish.AccountData
	var found bool
	var amap map[string]*redfish.AccountData
	argParse := flag.NewFlagSet("get-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Get detailed information for user identified by name")
	var id = argParse.String("id", "", "Get detailed information for user identified by ID")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *name != "" && *id != "" {
		return errors.New("ERROR: Options -name and -id are mutually exclusive")
	}

	if *name == "" && *id == "" {
		return errors.New("ERROR: Required options -name or -id not found")
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
		amap, err = r.MapAccountsById()
	} else {
		amap, err = r.MapAccountsByName()
	}

	if err != nil {
		return err
	}

	if *id != "" {
		acc, found = amap[*id]
	} else {
		acc, found = amap[*name]
	}

	if found {
		// XXX: Allow for different output formats like JSON, YAML, ... ?
		if *id != "" {
			fmt.Println(" " + *acc.Id)
		} else {
			fmt.Println(" " + *acc.UserName)
		}

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

	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *name, r.Hostname)
		}
	}

	return nil
}
