package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
	"syscall"
)

func modifyUser(r redfish.Redfish, args []string) error {
	var acc redfish.AccountCreateData

	argParse := flag.NewFlagSet("modify-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account to modify")
	var rename = argParse.String("rename", "", "Rename account to new name")
	var role = argParse.String("role", "", "New role of user account")
	var password = argParse.String("password", "", "New password for user account")
	var passwordFile = argParse.String("password-file", "", "Read password from file")
	var askPassword = argParse.Bool("ask-password", false, "New password for user account, will be read from stdin")
	var enable = argParse.Bool("enable", false, "Enable account")
	var disable = argParse.Bool("disable", false, "Disable account")
	var lock = argParse.Bool("lock", false, "Lock account")
	var unlock = argParse.Bool("unlock", false, "Unlock account")
	var hpePrivileges = argParse.String("hpe-privileges", "", "List of privileges for HP(E) systems")
	var err error

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *enable && *disable {
		return errors.New("ERROR: -enable and -disable are mutually exclusive")
	}

	if *password != "" && *passwordFile != "" {
		return errors.New("ERROR: -password and -password-file are mutually exclusive")
	}

	if (*password != "" || *passwordFile != "") && *askPassword {
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

	if *password != "" && *askPassword {
		return errors.New("ERROR: -password and -ask-password are mutually exclusive")
	}

	// Initialize session
	err = r.Initialise()
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

	if *rename != "" {
		acc.UserName = *name
	}

	// HP don't use or supports roles but their own privilege map
	if r.Flavor == redfish.RedfishHP {
		if *hpePrivileges != "" {
			acc.HPEPrivileges, err = hpeParsePrivileges(*hpePrivileges)
			if err != nil {
				return err
			}
		}
	} else {

		if *hpePrivileges != "" {
			log.WithFields(log.Fields{
				"hostname":      r.Hostname,
				"port":          r.Port,
				"timeout":       r.Timeout,
				"insecure_ssl":  r.InsecureSSL,
				"flavor":        r.Flavor,
				"flavor_string": r.FlavorString,
			}).Warning("This is not a HP(E) system, ignoring -hpe-privileges")
		}
	}
	if *role != "" {
		acc.Role = *role
	}

	// ask for password ?
	if *askPassword {
		fmt.Printf("Password for %s: ", *name)
		rawPass, _ := terminal.ReadPassword(int(syscall.Stdin))
		pass1 := strings.Replace(strings.Replace(strings.Replace(string(rawPass), "\r", "", -1), "\n", "", -1), "\t", "", -1)
		fmt.Println()

		fmt.Printf("Repeat password for %s: ", *name)
		rawPass, _ = terminal.ReadPassword(int(syscall.Stdin))
		pass2 := strings.Replace(strings.Replace(strings.Replace(string(rawPass), "\r", "", -1), "\n", "", -1), "\t", "", -1)
		fmt.Println()

		if pass1 != pass2 {
			return fmt.Errorf("ERROR: Passwords does not match for user %s", *name)
		}

		acc.Password = pass1
	}

	if *password != "" {
		acc.Password = *password
	} else if *passwordFile != "" {
		passwd, err := readSingleLine(*passwordFile)
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
