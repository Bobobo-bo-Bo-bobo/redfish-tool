package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

func Passwd(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("passwd", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account")
	var password = argParse.String("password", "", "Password for new user account. If omitted the password will be asked and read from stdin")

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

	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	// ask for password ?
	if *password == "" {
		fmt.Printf("Password for %s: ", *name)
		raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
		*password = strings.TrimSpace(string(raw_pass))
	}

	err = r.ChangePassword(*name, *password)
	if err != nil {
		return err
	}

	return nil
}
