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

func printManagerJson(r redfish.Redfish, mgr *redfish.ManagerData) string {
	var result string

	str, err := json.Marshal(mgr)
	if err != nil {
		log.Panic(err)
	}
	result = fmt.Sprintf("{\"%s\":%s}", r.Hostname, string(str))

	return result
}

func printManagerText(r redfish.Redfish, mgr *redfish.ManagerData) string {
	var result string

	if mgr.Id != nil {
		result += " Id: " + *mgr.Id + "\n"
	}
	if mgr.Name != nil {
		result += " Name:" + *mgr.Name + "\n"
	}

	if mgr.ManagerType != nil {
		result += " ManagerType:" + *mgr.ManagerType + "\n"
	}

	if mgr.UUID != nil {
		result += " UUID:" + *mgr.UUID + "\n"
	}

	if mgr.FirmwareVersion != nil {
		result += " FirmwareVersion:" + *mgr.FirmwareVersion + "\n"
	}

	result += " Status: " + "\n"
	if mgr.Status.State != nil {
		result += "  State: " + *mgr.Status.State + "\n"
	}
	if mgr.Status.Health != nil {
		result += "  Health: " + *mgr.Status.Health + "\n"
	}

	if mgr.SelfEndpoint != nil {
		result += " Endpoint: " + *mgr.SelfEndpoint + "\n"
	}

	return result
}

func printManager(r redfish.Redfish, mgr *redfish.ManagerData, format uint) string {
	if format == OUTPUT_JSON {
		return printManagerJson(r, mgr)
	}

	return printManagerText(r, mgr)
}

func GetManager(r redfish.Redfish, args []string, format uint) error {
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
		fmt.Println(printManager(r, mgr, format))
	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "User %s not found on %s\n", *uuid, r.Hostname)
		}
	}

	return nil
}
