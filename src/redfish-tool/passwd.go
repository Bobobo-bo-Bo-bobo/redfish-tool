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

func passwd(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("passwd", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account")
	var password = argParse.String("password", "", "Password for new user account. If omitted the password will be asked and read from stdin")
	var passwordFile = argParse.String("password-file", "", "Read password from file")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *name == "" {
		return errors.New("ERROR: Required options -name not found")
	}

	if *password != "" && *passwordFile != "" {
		return fmt.Errorf("ERROR: -password and -password-file are mutually exclusive")
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

	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	// ask for password ?
	if *password == "" {
		if *passwordFile == "" {
			fmt.Printf("Password for %s: ", *name)
			rawPass, _ := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			pass1 := strings.TrimSpace(string(rawPass))

			fmt.Printf("Repeat password for %s: ", *name)
			rawPass, _ = terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			pass2 := strings.TrimSpace(string(rawPass))

			if pass1 != pass2 {
				return fmt.Errorf("ERROR: Passwords does not match for user %s", *name)
			}

			*password = pass1
		} else {
			passwd, err := readSingleLine(*passwordFile)
			if err != nil {
				return err
			}

			password = &passwd
		}
	}

	err = r.ChangePassword(*name, *password)
	if err != nil {
		return err
	}

	return nil
}
