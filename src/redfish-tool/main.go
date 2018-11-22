package main

import (
	"fmt"
	"os"
	"redfish"
)

func main() {
	rf := redfish.Redfish{}
	rcfg := &redfish.RedfishConfiguration{
		Hostname: "localhost",
		Port:     8000,
		Username: "redfish",
		Password: "1t's s0 f1uffy I'm g0nn4 D13!",
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
