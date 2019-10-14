package main

import (
	"encoding/json"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
)

func printAllSystemsText(r redfish.Redfish, smap map[string]*redfish.SystemData) string {
	var result string

	result = r.Hostname + "\n"
	for sname, sdata := range smap {
		result += " " + sname + "\n"

		if sdata.ID != nil {
			result += "  Id:" + *sdata.ID + "\n"
		}

		if sdata.UUID != nil {
			result += "  UUID:" + *sdata.UUID + "\n"
		}

		if sdata.Name != nil {
			result += "  Name:" + *sdata.Name + "\n"
		}

		if sdata.SerialNumber != nil {
			result += "  SerialNumber:" + *sdata.SerialNumber + "\n"
		}

		if sdata.Manufacturer != nil {
			result += "  Manufacturer:" + *sdata.Manufacturer + "\n"
		}

		if sdata.Model != nil {
			result += "  Model:" + *sdata.Model + "\n"
		}

		result += "  Status:" + "\n"
		if sdata.Status.State != nil {
			result += "   State: " + *sdata.Status.State + "\n"
		}
		if sdata.Status.Health != nil {
			result += "   Health: " + *sdata.Status.Health + "\n"
		}
		if sdata.Status.HealthRollUp != nil {
			result += "   HealthRollUp: " + *sdata.Status.HealthRollUp + "\n"
		}

		if sdata.PowerState != nil {
			result += "  PowerState:" + *sdata.PowerState + "\n"
		}

		if sdata.BIOSVersion != nil {
			result += "  BIOSVersion:" + *sdata.BIOSVersion + "\n"
		}

		if sdata.SelfEndpoint != nil {
			result += "  SelfEndpoint:" + *sdata.SelfEndpoint + "\n"
		}

	}
	return result
}

func printAllSystemsJSON(r redfish.Redfish, smap map[string]*redfish.SystemData) string {
	var result string

	for _, sdata := range smap {
		str, err := json.Marshal(sdata)

		// Should NEVER happen!
		if err != nil {
			log.Panic(err)
		}

		result += fmt.Sprintf("{\"%s\":%s}\n", r.Hostname, string(str))
	}

	return result
}

func printAllSystems(r redfish.Redfish, smap map[string]*redfish.SystemData, format uint) string {
	if format == OutputJSON {
		return printAllSystemsJSON(r, smap)
	}
	return printAllSystemsText(r, smap)
}

func getAllSystems(r redfish.Redfish, format uint) error {
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

	// get all systems
	smap, err := r.MapSystemsByID()
	if err != nil {
		return err
	}

	fmt.Println(printAllSystems(r, smap, format))

	return nil
}
