package main

import (
	"encoding/json"
	"errors"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
)

func printAllRolesJSON(r redfish.Redfish, rmap map[string]*redfish.RoleData) string {
	var result string

	for _, rle := range rmap {
		str, err := json.Marshal(rle)
		// Should NEVER happen!
		if err != nil {
			log.Panic(err)
		}

		result += fmt.Sprintf("{\"%s\":%s}\n", r.Hostname, string(str))
	}

	return result
}

func printAllRolesText(r redfish.Redfish, rmap map[string]*redfish.RoleData) string {
	var result string

	result = r.Hostname + "\n"

	// loop over all endpoints
	for rid, rle := range rmap {
		result += " " + rid + "\n"
		if rle.ID != nil && *rle.ID != "" {
			result += "  Id: " + *rle.ID + "\n"
		}

		if rle.Name != nil && *rle.Name != "" {
			result += "  Name: " + *rle.Name + "\n"
		}

		if rle.IsPredefined != nil {
			if *rle.IsPredefined {
				result += "  IsPredefined: true" + "\n"
			} else {
				result += "  IsPredefined: false" + "\n"
			}
		}

		if len(rle.AssignedPrivileges) != 0 {
			result += "  Assigned privieleges" + "\n"
			for _, p := range rle.AssignedPrivileges {
				result += "   " + p + "\n"
			}
		}

		if rle.SelfEndpoint != nil {
			result += "  Endpoint: " + *rle.SelfEndpoint + "\n"
		}
	}
	return result
}

func printAllRoles(r redfish.Redfish, rmap map[string]*redfish.RoleData, format uint) string {
	if format == OutputJSON {
		return printAllRolesJSON(r, rmap)
	}

	return printAllRolesText(r, rmap)
}

func getAllRoles(r redfish.Redfish, format uint) error {
	// Initialize session
	err := r.Initialise()
	if err != nil {
		return fmt.Errorf("ERROR: Initialisation failed for %s: %s", r.Hostname, err.Error())
	}

	// Login
	err = r.Login()
	if err != nil {
		return fmt.Errorf("ERROR: Login to %s failed: %s", r.Hostname, err.Error())
	}

	defer r.Logout()

	// check if vendor support roles
	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HasAccountRoles != redfish.HasAccountRoles {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support roles")
		}
	}

	// get all role endpoints - Note: role names are _NOT_ unique but IDs are!
	rmap, err := r.MapRolesByID()
	if err != nil {
		return err
	}

	fmt.Println(printAllRoles(r, rmap, format))

	return nil
}
