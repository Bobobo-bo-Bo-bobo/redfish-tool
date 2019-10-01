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

func printSystemJson(r redfish.Redfish, sys *redfish.SystemData) string {
	var result string

	str, err := json.Marshal(sys)
	// Should NEVER happen!
	if err != nil {
		log.Panic(err)
	}

	result += fmt.Sprintf("{\"%s\":%s}\n", r.Hostname, string(str))

	return result
}

func printSystemText(r redfish.Redfish, sys *redfish.SystemData) string {
	var result string

	result = r.Hostname + "\n"

	if sys.Id != nil {
		result += "  Id:" + *sys.Id + "\n"
	}

	if sys.UUID != nil {
		result += "  UUID:" + *sys.UUID + "\n"
	}

	if sys.Name != nil {
		result += "  Name:" + *sys.Name + "\n"
	}

	if sys.SerialNumber != nil {
		result += "  SerialNumber:" + *sys.SerialNumber + "\n"
	}

	if sys.Manufacturer != nil {
		result += "  Manufacturer:" + *sys.Manufacturer + "\n"
	}

	if sys.Model != nil {
		result += "  Model:" + *sys.Model + "\n"
	}

	result += "  Status:" + "\n"
	if sys.Status.State != nil {
		result += "   State: " + *sys.Status.State + "\n"
	}
	if sys.Status.Health != nil {
		result += "   Health: " + *sys.Status.Health + "\n"
	}
	if sys.Status.HealthRollUp != nil {
		result += "   HealthRollUp: " + *sys.Status.HealthRollUp + "\n"
	}

	if sys.PowerState != nil {
		result += "  PowerState:" + *sys.PowerState + "\n"
	}

	if sys.BIOSVersion != nil {
		result += "  BIOSVersion:" + *sys.BIOSVersion + "\n"
	}

	if sys.SelfEndpoint != nil {
		result += "  SelfEndpoint:" + *sys.SelfEndpoint + "\n"
	}

	return result
}

func printSystem(r redfish.Redfish, sys *redfish.SystemData, format uint) string {
	if format == OUTPUT_JSON {
		return printSystemJson(r, sys)
	}

	return printSystemText(r, sys)
}

func GetSystem(r redfish.Redfish, args []string, format uint) error {
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

	if found {
		fmt.Println(r, sys, format)
	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "System %s not found on %s\n", *id, r.Hostname)
		} else {
			fmt.Fprintf(os.Stderr, "System %s not found on %s\n", *uuid, r.Hostname)
		}

	}
	return nil
}
