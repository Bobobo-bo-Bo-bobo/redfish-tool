package main

import (
	"errors"
	"flag"
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
)

func compareAndSetCSRField(s *string, a *string) *string {
	var _s string
	var _a string

	if s != nil {
		_s = *s
	}

	if a != nil {
		_a = *a
	}

	if _s == "" && _a != "" {
		return a
	} else if _s != "" && _a == "" {
		return s
	} else if _s != "" && _a != "" {
		return s
	} else {
		// a == "" && s == ""
		return s
	}
}

func GenCSR(r redfish.Redfish, args []string) error {
	var csrdata redfish.CSRData

	argParse := flag.NewFlagSet("gen-csr", flag.ExitOnError)

	var c = argParse.String("country", "", "CSR - country")
	var _c = argParse.String("c", "", "CSR - country")
	var s = argParse.String("state", "", "CSR - state or province")
	var _s = argParse.String("s", "", "CSR - state or province")
	var l = argParse.String("locality", "", "CSR - locality or city")
	var _l = argParse.String("l", "", "CSR - locality or city")
	var o = argParse.String("organisation", "", "CSR - organisation")
	var _o = argParse.String("o", "", "CSR - organisation")
	var ou = argParse.String("organisational-unit", "", "CSR - organisational unit")
	var _ou = argParse.String("ou", "", "CSR - organisational unit")
	var cn = argParse.String("common-name", "", "CSR - common name")
	var _cn = argParse.String("cn", "", "CSR - common name")

	argParse.Parse(args)

	c = compareAndSetCSRField(c, _c)
	s = compareAndSetCSRField(s, _s)
	l = compareAndSetCSRField(l, _l)
	o = compareAndSetCSRField(o, _o)
	ou = compareAndSetCSRField(ou, _ou)
	cn = compareAndSetCSRField(cn, _cn)

	// at least the common-name (CN) must be set, see Issue#3
	if *cn == "" {
		return errors.New(fmt.Sprintf("ERROR: At least the common name must be set for CSR generation"))
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

	// check if vendor support roles
	err = r.GetVendorFlavor()
	if err != nil {
		return err
	}

	capa, found := redfish.VendorCapabilities[r.FlavorString]
	if found {
		if capa&redfish.HAS_SECURITYSERVICE != redfish.HAS_SECURITYSERVICE {
			fmt.Println(r.Hostname)
			return errors.New("Vendor does not support CSR generation")
		}
	}

	csrdata = redfish.CSRData{
		C:  *c,
		S:  *s,
		L:  *l,
		O:  *o,
		OU: *ou,
		CN: *cn,
	}

	err = r.GenCSR(csrdata)
	if err != nil {
		return err
	}

	return nil
}
