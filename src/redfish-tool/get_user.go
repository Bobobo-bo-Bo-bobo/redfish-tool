package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
	"os"
)

func printUserText(r redfish.Redfish, acc *redfish.AccountData) string {
	var result string

	result = r.Hostname + "\n"

	if acc.ID != nil && *acc.ID != "" {
		result += " Id: " + *acc.ID + "\n"
	}

	if acc.Name != nil && *acc.Name != "" {
		result += " Name: " + *acc.Name + "\n"
	}

	if acc.UserName != nil && *acc.UserName != "" {
		result += " UserName: " + *acc.UserName + "\n"
	}

	if acc.Password != nil && *acc.Password != "" {
		result += " Password: " + *acc.Password + "\n"
	}

	if acc.RoleID != nil && *acc.RoleID != "" {
		result += " RoleId: " + *acc.RoleID + "\n"
	}

	if acc.Enabled != nil {
		if *acc.Enabled {
			result += " Enabled: true" + "\n"
		} else {
			result += " Enabled: false" + "\n"
		}
	}

	if acc.Locked != nil {
		if *acc.Locked {
			result += " Locked: true" + "\n"
		} else {
			result += " Locked: false" + "\n"
		}
	}

	if acc.SelfEndpoint != nil {
		result += " Endpoint: " + *acc.SelfEndpoint + "\n"
	}

	return result
}

func printUserJSON(r redfish.Redfish, acc *redfish.AccountData) string {
	var result string

	str, err := json.Marshal(acc)
	if err != nil {
		log.Panic(err)
	}
	result = fmt.Sprintf("{\"%s\":%s}", r.Hostname, string(str))

	return result
}

func printUser(r redfish.Redfish, acc *redfish.AccountData, format uint) string {
	if format == OutputJSON {
		return printUserJSON(r, acc)
	}

	return printUserText(r, acc)
}

func getUser(r redfish.Redfish, args []string, format uint) error {
	var acc *redfish.AccountData
	var found bool
	var amap map[string]*redfish.AccountData
	argParse := flag.NewFlagSet("get-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Get detailed information for user identified by name")
	var id = argParse.String("id", "", "Get detailed information for user identified by ID")

	argParse.Parse(args)

	if *name != "" && *id != "" {
		return errors.New("ERROR: Options -name and -id are mutually exclusive")
	}

	if *name == "" && *id == "" {
		return errors.New("ERROR: Required options -name or -id not found")
	}

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

	// get all account endpoints
	if *id != "" {
		amap, err = r.MapAccountsByID()
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
		fmt.Println(printUser(r, acc, format))
	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *name, r.Hostname)
		}
	}

	return nil
}
