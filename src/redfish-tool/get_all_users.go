package main

import (
	"fmt"

	"redfish"
)

func GetAllUsers(r redfish.Redfish, cfg *redfish.RedfishConfiguration) error {
	// get all account endpoints
	ael, err := r.GetAccounts(cfg)
	if err != nil {
		return err
	}

	fmt.Println(cfg.Hostname)
	// loop over all endpoints
	for _, ae := range ael {
		acc, err := r.GetAccountData(cfg, ae)
		if err != nil {
			return err
		}

		// XXX: Allow for different output formats like JSON, YAML, ... ?
		fmt.Println(" " + *acc.UserName)
	}

	return err
}
