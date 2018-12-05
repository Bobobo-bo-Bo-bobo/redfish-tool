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

func AddUser(r redfish.Redfish, args []string) error {
	var acc redfish.AccountCreateData

	argParse := flag.NewFlagSet("add-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account to create")
	var role = argParse.String("role", "", "Role of user account to create")
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

	acc.UserName = *name

	// HP don't use or supports roles but their own privilege map
	if r.Flavor == redfish.REDFISH_HP {
		// FIXME: handle this!
	} else {
		if *role == "" {
			return errors.New("ERROR: Required option -role not found")
		}
		acc.Role = *role
	}

	// ask for password ?
	if *password == "" {
        fmt.Print("Password for " + *name + ": ")
		raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
		acc.Password = strings.TrimSpace(string(raw_pass))
	} else {
        acc.Password = *password
    }

	err = r.AddAccount(acc)
	if err != nil {
		return err
	}

	return nil
}
