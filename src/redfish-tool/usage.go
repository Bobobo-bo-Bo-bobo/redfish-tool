package main

import (
	"fmt"
)

func ShowUsage() {
	const version string = "1.0.0"

	fmt.Printf("redfish-tool version %s\n"+
		"Copyright (C) 2018 - 2019 by Andreas Maus <maus@ypbind.de>\n"+
		"This program comes with ABSOLUTELY NO WARRANTY.\n"+
		"\n"+
		"redfish-tool is distributed under the Terms of the GNU General\n"+
		"Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)\n"+
		"\n"+
		"Usage redfish-tool [-ask] [-config=<cfg>] [-help] [-password=<pass>] [-user=<user>] <command> [<cmd_options>]\n"+
		"\n"+
		"Global options:\n"+
		"\n"+
		"  -ask\n"+
		"    	Ask for password\n"+
		"  -config=<cfg>\n"+
		"    	Configuration file to use\n"+
		"  -help\n"+
		"    	Show help text\n"+
		"  -insecure\n"+
		"    	Skip SSL certificate verification\n"+
		"  -password=<pass>\n"+
		"    	Password to use for authentication\n"+
		"  -user=<user>\n"+
		"    	Username to use for authentication\n"+
		"  -verbose\n"+
		"    	Verbose operation\n"+
		"\n"+
		"Commands:\n"+
		"\n"+
		"  get-all-users\n"+
		"\n"+
		"  get-user\n"+
		"    -name=<name>\n"+
		"         Get detailed information for user identified by name (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for user identified by ID (*)\n"+
		"\n"+
		"    (*) -name and -id are mutually exclusive\n"+
		"\n"+
		"  get-all-roles\n"+
		"\n"+
		"  get-role\n"+
		"    -id=<id>\n"+
		"         Get detailed information for role identified by ID (**)\n"+
		"\n"+
		"    (**) Because role names are not unique, roles can only be listed by ID\n"+
		"\n"+
		"  get-all-managers\n"+
		"\n"+
		"  get-manager\n"+
		"    -uuid=<uuid>\n"+
		"         Get detailed information for user identified by UUID (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for user identified by ID (*)\n"+
		"\n"+
		"    (*) -uuid and -id are mutually exclusive\n"+
		"\n"+
		"  get-all-systems\n"+
		"\n"+
		"  get-system\n"+
		"    -uuid=<uuid>\n"+
		"         Get detailed information for system identified by UUID (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for system identified by ID (*)\n"+
		"\n"+
		"    (*) -uuid and -id are mutually exclusive\n"+
		"\n"+
		"  gen-csr\n"+
		"    -country=<c>\n"+
		"       CSR - country\n"+
		"    -state=<s>\n"+
		"       CSR - state or province\n"+
		"    -locality=<l>\n"+
		"       CSR - locality or city\n"+
		"   -organisation=<o>\n"+
		"       CSR - organisation\n"+
		"   -organisational-unit=<ou>\n"+
		"       CSR - organisational unit\n"+
		"   -common-name=<cn>\n"+
		"       CSR - common name, hostname will be used if no CN is set\n"+
		"\n"+
		"\n", version)
}
