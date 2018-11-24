package main

import (
	"bufio"
	"fmt"
	"os"
	"redfish"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	rf := redfish.Redfish{}
	// Testing with local simulator:
	//
	// https://github.com/DMTF/Redfish-Profile-Simulator
	//
	// User: root
	// Password: password123456
	// Returned token from SessionService: 123456SESSIONauthcode
	//
	/*
		rcfg := &redfish.RedfishConfiguration{
			Hostname: "localhost",
			Port:     5000,
			Username: "root",
			Password: "password123456",
		}
	*/
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Hostname: ")
	host, _ := r.ReadString('\n')
	host = strings.TrimSpace(host)

	fmt.Print("User: ")
	user, _ := r.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Password: ")
	raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
	rcfg := &redfish.RedfishConfiguration{
		Hostname:    host,
		Username:    user,
		Password:    strings.TrimSpace(string(raw_pass)),
		InsecureSSL: true,
	}

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
