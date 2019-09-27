package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"io/ioutil"
	"os"
)

func AddLicense(r redfish.Redfish, args []string) error {
	argParse := flag.NewFlagSet("add-license", flag.ExitOnError)
	var id = argParse.String("id", "", "Management board identified by ID")
	var uuid = argParse.String("uuid", "", "Management board identified by UUID")
	var l = argParse.String("license", "", "License data to add")
	var lf = argParse.String("license-file", "", "License file containing the license data")
	var mmap map[string]*redfish.ManagerData
	var mgr *redfish.ManagerData
	var ldata []byte
	var err error

	argParse.Parse(args)

	fmt.Println(r.Hostname)

	if *uuid != "" && *id != "" {
		return errors.New("ERROR: Options -uuid and -id are mutually exclusive")
	}
	if *uuid == "" && *id == "" {
		return errors.New("ERROR: Required options -uuid or -id not found")
	}

	if *l != "" && *lf != "" {
		return errors.New("ERROR: Options -license and -license-file are mutually exclusive")
	}
	if *l == "" && *lf == "" {
		return errors.New("ERROR: Mandatory options -license or -license-file are not found")
	}

	if *lf != "" {
		if *lf == "-" {
			ldata, err = ioutil.ReadAll(os.Stdin)
		} else {
			ldata, err = ioutil.ReadFile(*lf)
		}
		if err != nil {
			return err
		}
	}
	if *l != "" {
		ldata = []byte(*l)
	}

	// Initialize session
	err = r.Initialise()
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
		err := r.AddLicense(mgr, ldata)
		if err != nil {
			return nil
		}
	} else {
		if *id != "" {
			fmt.Fprintf(os.Stderr, "Manager with ID %s not found on %s\n", *id, r.Hostname)
		} else if *uuid != "" {
			fmt.Fprintf(os.Stderr, "Manager with UUID %s not found on %s\n", *uuid, r.Hostname)
		}
	}

	return nil
}
