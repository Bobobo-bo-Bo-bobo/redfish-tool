package main

import (
	"fmt"
	"os"
	"redfish"
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
	rcfg := &redfish.RedfishConfiguration{
		Hostname: "localhost",
		Port:     5000,
		Username: "root",
		Password: "password123456",
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

    fmt.Println("---")
    fmt.Printf("%+v\n", rcfg)
    fmt.Println("---")
	os.Exit(0)
}
