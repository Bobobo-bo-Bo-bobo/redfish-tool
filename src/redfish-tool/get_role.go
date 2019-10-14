package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
	"os"
)

func printRoleJSON(r redfish.Redfish, rle *redfish.RoleData) string {
	var result string

	str, err := json.Marshal(rle)
	if err != nil {
		log.Panic(err)
	}
	result = fmt.Sprintf("{\"%s\":%s}", r.Hostname, string(str))

	return result

}

func printRoleText(r redfish.Redfish, rle *redfish.RoleData) string {
	var result string

	result = r.Hostname + "\n"
	if rle.ID != nil && *rle.ID != "" {
		result += " Id: " + *rle.ID + "\n"
	}

	if rle.Name != nil && *rle.Name != "" {
		result += " Name: " + *rle.Name + "\n"
	}

	if rle.IsPredefined != nil {
		if *rle.IsPredefined {
			result += " IsPredefined: true" + "\n"
		} else {
			result += " IsPredefined: false" + "\n"
		}
	}

	if len(rle.AssignedPrivileges) != 0 {
		result += " Assigned privieleges" + "\n"
		for _, p := range rle.AssignedPrivileges {
			result += "   " + p + "\n"
		}
	}

	if rle.SelfEndpoint != nil {
		result += " Endpoint: " + *rle.SelfEndpoint + "\n"
	}
	return result
}

func printRole(r redfish.Redfish, rle *redfish.RoleData, format uint) string {
	if format == OutputJSON {
		return printRoleJSON(r, rle)
	}

	return printRoleText(r, rle)
}

func getRole(r redfish.Redfish, args []string, format uint) error {
	var rle *redfish.RoleData
	var found bool
	var rmap map[string]*redfish.RoleData
	argParse := flag.NewFlagSet("get-role", flag.ExitOnError)

	var id = argParse.String("id", "", "Get detailed information for role identified by ID")

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *id == "" {
		return errors.New("ERROR: Required option -id not found")
	}

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

	// get all roles
	rmap, err = r.MapRolesByID()

	if err != nil {
		return err
	}

	rle, found = rmap[*id]

	if found {
		fmt.Println(printRole(r, rle, format))
	} else {
		fmt.Fprintf(os.Stderr, "Role %s not found on %s\n", *id, r.Hostname)
	}

	return nil
}
