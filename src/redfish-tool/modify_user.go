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

func ModifyUser(r redfish.Redfish, args []string) error {
	var acc redfish.AccountCreateData

	argParse := flag.NewFlagSet("modify-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account to modify")
	var rename = argParse.String("rename", "", "Rename account to new name")
	var role = argParse.String("role", "", "New role of user account")
	var password = argParse.String("password", "", "New password for user account")
	var password_file = argParse.String("password-file", "", "Read password from file")
	var ask_password = argParse.Bool("ask-password", false, "New password for user account, will be read from stdin")
	var enable = argParse.Bool("enable", false, "Enable account")
	var disable = argParse.Bool("disable", false, "Disable account")
	var lock = argParse.Bool("lock", false, "Lock account")
	var unlock = argParse.Bool("unlock", false, "Unlock account")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *enable && *disable {
		return errors.New("ERROR: -enable and -disable are mutually exclusive")
	}

	if *password != "" && *password_file != "" {
		return errors.New("ERROR: -password and -password-file are mutually exclusive")
	}

	if (*password != "" || *password_file != "") && *ask_password {
		return errors.New("ERROR: -password/-password-file and -ask-password are mutually exclusive")
	}

	if *enable {
		acc.Enabled = enable
	}

	if *disable {
		e := false
		acc.Enabled = &e
	}

	if *lock && *unlock {
		return errors.New("ERROR: -lock and -unlock are mutually exclusive")
	}

	if *lock {
		acc.Locked = lock
	}

	if *unlock {
		l := false
		acc.Locked = &l
	}

	if *name == "" {
		return errors.New("ERROR: Required options -name not found")
	}

	if *password != "" && *ask_password {
		return errors.New("ERROR: -password and -ask-password are mutually exclusive")
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

	if *rename != "" {
		acc.UserName = *name
	}

	// HP don't use or supports roles but their own privilege map
	if r.Flavor == redfish.REDFISH_HP {
		// FIXME: handle this!
	} else {
		if *role != "" {
			acc.Role = *role
		}
	}

	// ask for password ?
	if *ask_password {
		fmt.Printf("Password for %s: ", *name)
		raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
		pass1 := strings.TrimSpace(string(raw_pass))
		fmt.Println()

		fmt.Printf("Repeat password for %s: ", *name)
		raw_pass, _ = terminal.ReadPassword(int(syscall.Stdin))
		pass2 := strings.TrimSpace(string(raw_pass))
		fmt.Println()

		if pass1 != pass2 {
			return errors.New(fmt.Sprintf("ERROR: Passwords does not match for user %s", *name))
		}

		acc.Password = pass1
	}

	if *password != "" {
		acc.Password = *password
	} else if *password_file != "" {
		passwd, err := ReadSingleLine(*password_file)
		if err != nil {
			return err
		}
		acc.Password = passwd
	}

	err = r.ModifyAccount(*name, acc)
	if err != nil {
		return err
	}

	return nil
}
