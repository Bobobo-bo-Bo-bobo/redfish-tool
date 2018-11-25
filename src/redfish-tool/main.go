package main

import (
	"flag"
	"fmt"
	"os"
	"redfish"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	insecure := flag.Bool("insecure", false, "Skip SSL certificate verification")
	verbose := flag.Bool("verbose", false, "Verbose operation")
	ask := flag.Bool("ask", false, "Ask for password")
	user := flag.String("user", "", "Username to use for authentication")
	password := flag.String("password", "", "Password to use for authentication")
	config_file := flag.String("config", "", "Configuration file to use")
	help := flag.Bool("help", false, "Show help text")
	hosts := flag.String("host", "", "Hosts to work on")
	port := flag.Int("port", 0, "Alternate port to connect to")
	timeout := flag.Int64("timeout", 60, "Connection timeout in seconds")

	flag.Usage = ShowUsage
	flag.Parse()
	if *help {
		ShowUsage()
		os.Exit(0)
	}

	trailing := flag.Args()

	if *config_file != "" {
		// read and parse configuration file
	} else {
		if *ask {
			fmt.Print("Password: ")
			raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
			pass := strings.TrimSpace(string(raw_pass))
			password = &pass
			fmt.Println()
		}
	}

	// get requested command
	if len(trailing) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command defined\n\n")
		ShowUsage()
		os.Exit(1)
	}
	command := strings.ToLower(trailing[0])

	if *hosts == "" {
		fmt.Fprintf(os.Stderr, "Error: No destination host given\n\n")
		ShowUsage()
		os.Exit(1)
	}

	if *timeout < 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid timeout %d; must be >= 0\n\n", *timeout)
		os.Exit(2)
	}

	host_list := strings.Split(*hosts, ",")
	for _, host := range host_list {
		rcfg := &redfish.RedfishConfiguration{
			Hostname:    host,
			Port:        *port,
			Username:    *user,
			Password:    *password,
			InsecureSSL: *insecure,
			Verbose:     *verbose,
			Timeout:     time.Duration(*timeout) * time.Second,
		}

		// XXX: optionally parse additional and command specific command lines

		rf := redfish.Redfish{}

		// Initialize session
		err = rf.Initialise(rcfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Initialisation failed for %s: %s\n", host, err.Error())
			continue
		}

		// Login
		err = rf.Login(rcfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Login to %s failed: %s\n", host, err.Error())
			continue
		}

		if command == "get-all-users" {
			err = GetAllUsers(rf, rcfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		} // XXX: ...

		// Logout
		err = rf.Logout(rcfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Logout from %s failed: %s\n", host, err.Error())
		}
	}

	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}

	/*
		rcfg := &redfish.RedfishConfiguration{
			Hostname:    hostname,
			Username:    *user,
			Password:    *password,
			InsecureSSL: *insecure,
			Verbose:     *verbose,
		}

		fmt.Println("")
		fmt.Print("Initialise - ")
		err := rf.Initialise(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
		}

		fmt.Print("Login - ")
		err = rf.Login(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf(" + Session stored at %s\n", *rcfg.SessionLocation)
			fmt.Printf(" + X-Auth-Token: %s\n", *rcfg.AuthToken)
		}

		fmt.Print("Systems - ")
		sys, err := rf.GetSystems(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf(" + %d systems reported\n", len(sys))

			for _, s := range sys {
				fmt.Printf("  * %s\n", s)
			}
		}

		fmt.Printf("System: %s - ", sys[0])
		ssys, err := rf.GetSystemData(rcfg, sys[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf("%+v\n", ssys)
		}

		fmt.Print("Accounts - ")
		accs, err := rf.GetAccounts(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf(" + %d accounts reported\n", len(accs))
			for _, a := range accs {
				fmt.Printf("  * %s\n", a)
			}
		}

		fmt.Printf("Account: %s - ", accs[0])
		acc, err := rf.GetAccountData(rcfg, accs[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf("%+v\n", acc)
		}

		fmt.Print("Roles - ")
		roles, err := rf.GetRoles(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf(" + %d roles reported\n", len(roles))
			for _, a := range roles {
				fmt.Printf("  * %s\n", a)
			}
		}

		fmt.Printf("Role: %s - ", roles[0])
		role, err := rf.GetRoleData(rcfg, roles[0])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
			fmt.Printf("%+v\n", role)
		}

		fmt.Print("Logout - ")
		err = rf.Logout(rcfg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("OK")
		}
	*/
}
