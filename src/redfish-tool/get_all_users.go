package main

import (
	"fmt"

	"redfish"
)

func GetAllUsers(r redfish.Redfish) error {
	// get all account endpoints
	ael, err := r.GetAccounts()
	if err != nil {
		return err
	}

	fmt.Println(r.Hostname)
	// loop over all endpoints
	for _, ae := range ael {
		acc, err := r.GetAccountData(ae)
		if err != nil {
			return err
		}

		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + *acc.UserName)
	}

	return err
}
