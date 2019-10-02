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

func printLicenseJson(r redfish.Redfish, l *redfish.ManagerLicenseData) string {
	var result string

	str, err := json.Marshal(l)
	if err != nil {
		log.Panic(err)
	}
	result = fmt.Sprintf("{\"%s\":%s}", r.Hostname, string(str))

	return result
}

func printLicenseText(r redfish.Redfish, l *redfish.ManagerLicenseData) string {
	var result string

	result = r.Hostname + "\n"

	if l.Name != "" {
		result += " Name: " + l.Name + "\n"
	} else {
		result += " Name: -" + "\n"
	}

	if l.Type != "" {
		result += " Type: " + l.Type + "\n"
	} else {
		result += " Type: -" + "\n"
	}

	if l.Expiration != "" {
		result += " Expiration: " + l.Expiration + "\n"
	} else {
		result += " Expiration: -" + "\n"
	}

	if l.License != "" {
		result += " License: " + l.License + "\n"
	} else {
		result += " License: -" + "\n"
	}

	return result
}

func printLicense(r redfish.Redfish, l *redfish.ManagerLicenseData, format uint) string {
	if format == OUTPUT_JSON {
		return printLicenseJson(r, l)
	}

	return printLicenseText(r, l)
}

func GetLicense(r redfish.Redfish, args []string, format uint) error {
	argParse := flag.NewFlagSet("get-license", flag.ExitOnError)
	var id = argParse.String("id", "", "Management board identified by ID")
	var uuid = argParse.String("uuid", "", "Management board identified by UUID")
	var mmap map[string]*redfish.ManagerData
	var mgr *redfish.ManagerData

	argParse.Parse(args)

	if *uuid != "" && *id != "" {
		return errors.New("ERROR: Options -uuid and -id are mutually exclusive")
	}

	if *uuid == "" && *id == "" {
		return errors.New("ERROR: Required options -uuid or -id not found")
	}

	// Initialize session
	err := r.Initialise()
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Initialisation failed for %s: %s\n", r.Hostname, err.Error()))
	}

	// Login
	err = r.Login()
	if err != nil {
		return errors.New(fmt.Sprintf("ERROR: Login to %s failed: %s\n", r.Hostname, err.Error()))
	}

	defer r.Logout()

	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HAS_LICENSE != redfish.HAS_LICENSE {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support license operations")
		}
	}

	if *id != "" {
		mmap, err = r.MapManagersById()
	} else if *uuid != "" {
		mmap, err = r.MapManagersByUuid()
	}

	if err != nil {
		return err
	}

	if *id != "" {
		mgr, found = mmap[*id]
	} else if *uuid != "" {
		mgr, found = mmap[*uuid]
	}

	if found {
		l, err := r.GetLicense(mgr)
		if err != nil {
			return err
		}

		fmt.Println(printLicense(r, l, format))

	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "Manager with ID %s not found on %s\n", *id, r.Hostname)
		} else if *uuid != "" {
			fmt.Fprintf(os.Stderr, "Manager with UUID %s not found on %s\n", *uuid, r.Hostname)
		}
	}

	return nil
}
