package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func DelUser(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("del-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account to remove")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *name == "" {
		return errors.New("ERROR: Required options -name not found")
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

	err = r.DeleteAccount(*name)
	if err != nil {
		return err
	}

	return nil
}
