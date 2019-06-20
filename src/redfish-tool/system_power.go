package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func SystemPower(r redfish.Redfish, args []string) error {
	var sys *redfish.SystemData
	var found bool
	var smap map[string]*redfish.SystemData

	argParse := flag.NewFlagSet("system-power", flag.ExitOnError)

	var uuid = argParse.String("uuid", "", "Get detailed information for system identified by UUID")
	var id = argParse.String("id", "", "Get detailed information for system identified by ID")
	var state = argParse.String("state", "", "Set power state of the system")

	argParse.Parse(args)

	if *uuid != "" && *id != "" {
		return errors.New("ERROR: Options -uuid and -id are mutually exclusive")
	}

	if *uuid == "" && *id == "" {
		return errors.New("ERROR: Required options -uuid or -id not found")
	}

	if *state == "" {
		return errors.New("ERROR: Option -state is mandatory")
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

	if !found {
		return errors.New("ERROR: Can't find system with requested ID/UUID")
	}

	err = r.SetSytemPowerState(sys, *state)
	return err
}
