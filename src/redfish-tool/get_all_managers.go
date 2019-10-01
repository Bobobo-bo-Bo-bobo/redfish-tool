package main

import (
	"encoding/json"
	"errors"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
)

func printAllManagersJson(r redfish.Redfish, mmap map[string]*redfish.ManagerData) string {
	var result string

	fmt.Println(r.Hostname)
	// loop over all endpoints
	for mname, mgr := range mmap {
		result += " " + mname + "\n"
		if mgr.Id != nil {
			result += "  Id: " + *mgr.Id + "\n"
		}
		if mgr.Name != nil {
			result += "  Name:" + *mgr.Name + "\n"
		}

		if mgr.ManagerType != nil {
			result += "  ManagerType:" + *mgr.ManagerType + "\n"
		}

		if mgr.UUID != nil {
			result += "  UUID:" + *mgr.UUID + "\n"
		}

		if mgr.FirmwareVersion != nil {
			result += "  FirmwareVersion:" + *mgr.FirmwareVersion + "\n"
		}

		result += "  Status: " + "\n"
		if mgr.Status.State != nil {
			result += "   State: " + *mgr.Status.State + "\n"
		}
		if mgr.Status.Health != nil {
			result += "   Health: " + *mgr.Status.Health + "\n"
		}
		if mgr.SelfEndpoint != nil {
			result += "  Endpoint: " + *mgr.SelfEndpoint + "\n"
		}

	}

	return result
}

func printAllManagersText(r redfish.Redfish, mmap map[string]*redfish.ManagerData) string {
	var result string

	for _, mgr := range mmap {
		str, err := json.Marshal(mgr)
		// Should NEVER happen!
		if err != nil {
			log.Panic(err)
		}

		result += fmt.Sprintf("{\"%s\":%s}\n", r.Hostname, string(str))
	}

	return result
}

func printAllManagers(r redfish.Redfish, mmap map[string]*redfish.ManagerData, format uint) string {
	if format == OUTPUT_JSON {
		return printAllManagersJson(r, mmap)
	}

	return printAllManagersText(r, mmap)
}

func GetAllManagers(r redfish.Redfish, format uint) error {
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

	fmt.Println(printAllManagers(r, mmap, format))

	return nil
}
