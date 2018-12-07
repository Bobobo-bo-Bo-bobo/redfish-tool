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
	var password_file = argParse.String("password-file", "", "Read password from file")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *name == "" {
		return errors.New("ERROR: Required options -name not found")
	}

	if *password != "" && *password_file != "" {
		return errors.New(fmt.Sprintf("ERROR: -password and -password-file are mutually exclusive"))
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
		if *password_file == "" {
			fmt.Printf("Password for %s: ", *name)
			raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			pass1 := strings.TrimSpace(string(raw_pass))

			fmt.Printf("Repeat password for %s: ", *name)
			raw_pass, _ = terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			pass2 := strings.TrimSpace(string(raw_pass))

			if pass1 != pass2 {
				return errors.New(fmt.Sprintf("ERROR: Passwords does not match for user %s", *name))
			}

			*password = pass1
		} else {
			passwd, err := ReadSingleLine(*password_file)
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
