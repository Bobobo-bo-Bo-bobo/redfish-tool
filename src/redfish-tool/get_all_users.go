package main

import (
	"encoding/json"
	"errors"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
)

func printAllUsersText(r redfish.Redfish, amap map[string]*redfish.AccountData) string {
	var result string

	result = r.Hostname + "\n"

	// loop over all endpoints
	for aname, acc := range amap {

		result += " " + aname + "\n"

		if acc.Id != nil && *acc.Id != "" {
			result += "  Id: " + *acc.Id + "\n"
		}

		if acc.Name != nil && *acc.Name != "" {
			result += "  Name: " + *acc.Name + "\n"
		}

		if acc.UserName != nil && *acc.UserName != "" {
			result += "  UserName: " + *acc.UserName + "\n"
		}

		if acc.Password != nil && *acc.Password != "" {
			result += "  Password: " + *acc.Password + "\n"
		}

		if acc.RoleId != nil && *acc.RoleId != "" {
			result += "  RoleId: " + *acc.RoleId + "\n"
		}

		if acc.Enabled != nil {
			if *acc.Enabled {
				result += "  Enabled: true" + "\n"
			} else {
				result += "  Enabled: false" + "\n"
			}
		}

		if acc.Locked != nil {
			if *acc.Locked {
				result += "  Locked: true" + "\n"
			} else {
				result += "  Locked: false" + "\n"
			}
		}

		if acc.SelfEndpoint != nil {
			result += "  Endpoint: " + *acc.SelfEndpoint + "\n"
		}
	}
	return result
}

func printAllUsersJson(r redfish.Redfish, amap map[string]*redfish.AccountData) string {
	var result string

	for _, acc := range amap {
		str, err := json.Marshal(acc)
		// Should NEVER happen!
		if err != nil {
			log.Panic(err)
		}

		result += fmt.Sprintf("{\"%s\":%s}\n", r.Hostname, string(str))
	}

	return result
}

func printAllUsers(r redfish.Redfish, amap map[string]*redfish.AccountData, format uint) string {
	if format == OUTPUT_JSON {
		return printAllUsersJson(r, amap)
	}
	return printAllUsersText(r, amap)
}

func GetAllUsers(r redfish.Redfish, format uint) error {
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
	amap, err := r.MapAccountsByName()
	if err != nil {
		return err
	}

	fmt.Println(printAllUsers(r, amap, format))
	return nil
}
