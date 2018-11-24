package main

import (
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

	rf := redfish.Redfish{}

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
	rcfg := &redfish.RedfishConfiguration{
		//	Hostname:    hostname,
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

	fmt.Print("Logout - ")
	err = rf.Logout(rcfg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
	os.Exit(0)
}
