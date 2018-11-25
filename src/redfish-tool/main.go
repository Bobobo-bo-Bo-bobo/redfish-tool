package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"redfish"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	insecure := flag.Bool("insecure", false, "Skip SSL certificate verification")
	verbose := flag.Bool("verbose", false, "Verbose operation")
	ask := flag.Bool("ask", false, "Ask for password")
	user := flag.String("user", "", "Username to use for authentication")
	password := flag.String("password", "", "Password to use for authentication")
	config_file := flag.String("config", "", "Configuration file to use")
	help := flag.Bool("help", false, "Show help text")

	flag.Usage = ShowUsage
	flag.Parse()
	if *help {
		ShowUsage()
		os.Exit(0)
	}

	trailing := flag.Args()

	rf := redfish.Redfish{}

	// FIXME: read data from command line/config file instead of asking for it
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Hostname: ")
	hostname, _ := r.ReadString('\n')
	hostname = strings.TrimSpace(hostname)

	if *config_file != "" {
		// read and parse configuration file
	} else {
		if *ask {
			fmt.Print("Password: ")
			raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
			pass := strings.TrimSpace(string(raw_pass))
			password = &pass
		}
	}

	// get requested command
	if len(trailing) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command defined\n\n")
		ShowUsage()
		os.Exit(1)
	}

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

	fmt.Print("Logout - ")
	err = rf.Logout(rcfg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
	os.Exit(0)
}
