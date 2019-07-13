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

func hpeParsePrivileges(privileges string) (uint, error) {
	var result uint = 0

	for _, priv := range strings.Split(strings.ToLower(privileges), ",") {
		_priv := strings.TrimSpace(priv)
		_bit, found := redfish.HPEPrivilegeMap[_priv]
		if !found {
			return 0, errors.New(fmt.Sprintf("ERROR: Unknown privilege %s", _priv))
		}
		result |= _bit
	}
	return result, nil
}

func AddUser(r redfish.Redfish, args []string) error {
	var acc redfish.AccountCreateData

	argParse := flag.NewFlagSet("add-user", flag.ExitOnError)

	var name = argParse.String("name", "", "Name of user account to create")
	var role = argParse.String("role", "", "Role of user account to create")
	var hpe_privileges = argParse.String("hpe-privileges", "", "List of privileges for HP(E) systems")
	var password = argParse.String("password", "", "Password for new user account. If omitted the password will be asked and read from stdin")
	var disabled = argParse.Bool("disabled", false, "Created account is disabled")
	var locked = argParse.Bool("locked", false, "Created account is locked")
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

	acc.UserName = *name

	// HP don't use or supports roles but their own privilege map
	if r.Flavor == redfish.REDFISH_HP {
		if *hpe_privileges != "" {
			acc.HPEPrivileges, err = hpeParsePrivileges(*hpe_privileges)
			if err != nil {
				return err
			}
		}
	} else {
		if *role == "" {
			return errors.New("ERROR: Required option -role not found")
		}
		if *hpe_privileges != "" {
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
	acc.Role = *role

	// ask for password ?
	if *password == "" {
		if *password_file == "" {
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
		} else {
			passwd, err := ReadSingleLine(*password_file)
			if err != nil {
				return err
			}
			acc.Password = passwd
		}
	} else {
		acc.Password = *password
	}

	if *disabled {
		enabled := false
		acc.Enabled = &enabled
	}

	if *locked {
		acc.Locked = locked
	}

	err = r.AddAccount(acc)
	if err != nil {
		return err
	}

	return nil
}
