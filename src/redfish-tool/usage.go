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
		"  get-all-users - List all users from service processor\n"+
		"\n"+
		"  get-user - List specific user from service processor\n"+
		"    -name=<name>\n"+
		"         Get detailed information for user identified by name (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for user identified by ID (*)\n"+
		"\n"+
		"    (*) -name and -id are mutually exclusive\n"+
		"\n"+
		"  get-all-roles - List all roles from service processor (*)\n"+
		"\n"+
		"    (*) HP/HPE servers do not defined roles\n"+
		"        Inspur servers report AccountService endpoint but don't implement it\n"+
		"\n"+
		"  get-role - List sepcific role (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for role identified by ID (**)\n"+
		"\n"+
		"    (*) HP/HPE servers do not defined roles\n"+
		"        Inspur servers report AccountService endpoint but don't implement it\n"+
		"\n"+
		"    (**) Because role names are not unique, roles can only be listed by ID\n"+
		"\n"+
		"  get-all-managers - List all managers\n"+
		"\n"+
		"  get-manager - List specific manager\n"+
		"    -uuid=<uuid>\n"+
		"         Get detailed information for user identified by UUID (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for user identified by ID (*)\n"+
		"\n"+
		"    (*) -uuid and -id are mutually exclusive\n"+
		"\n"+
		"  get-all-systems - List all systems\n"+
		"\n"+
		"  get-system - List specific system\n"+
		"    -uuid=<uuid>\n"+
		"         Get detailed information for system identified by UUID (*)\n"+
		"    -id=<id>\n"+
		"         Get detailed information for system identified by ID (*)\n"+
		"\n"+
		"    (*) -uuid and -id are mutually exclusive\n"+
		"\n"+
		"  gen-csr - Generate certificate signing request (*)\n"+
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
		"   (*) Inspur and Supermicro service processors do not provide endpoints for generation of certificate signing requests\n"+
		"\n"+
		"  fetch-csr - Fetch generated certificate signing request (*)\n"+
		"\n"+
		"   (*) Inspur and Supermicro service processors do not provide endpoints for generation of certificate signing requests\n"+
		"\n"+
		"  import-cert - Import certificate in PEM format (*)\n"+
		"    -certificate=<cert>\n"+
		"       Certificate file in PEM format to import\n"+
		"\n"+
		"   (*) Inspur and Supermicro service processors do not provide endpoints for certificate import\n"+
		"\n"+
		"  reset-sp - Reset service processor\n"+
		"\n"+
		"  add-user - Create a new user\n"+
		"    -name=<name>\n"+
		"        Name of user account to create\n"+
		"    -role=<role>\n"+
		"        Role of user account to create\n"+
		"    -password=<pass>\n"+
		"        Password for new user account. If omitted the password will be asked and read from stdin\n"+
		"    -password-file=<file>\n"+
		"        Read password from <file>. The password MUST be the first line in the file, all other lines are ignored\n"+
		"    -disabled\n"+
		"        Account is created but disabled\n"+
		"    -locked\n"+
		"        Account is created and locked\n"+
		"\n"+
		"  del-user - Delete an existing account\n"+
		"    -name=<name>\n"+
		"        Name of the account to delete\n"+
		"\n"+
		"  passwd - Change password of an existing account\n"+
		"    -name=<name>\n"+
		"        Name of the user account\n"+
		"    -password=<pass>\n"+
		"        New password. If omitted the password will be asked and read from stdin\n"+
		"    -password-file=<file>\n"+
		"        Read new password from <file>. The password MUST be the first line in the file, all other lines are ignored\n"+
		"\n"+
		"  modify-user - Modify an existing user\n"+
		"    -name=<name>\n"+
		"        Name of user account to modify\n"+
		"    -rename=<new_name>\n"+
		"        New name of user account\n"+
		"    -role=<role>\n"+
		"        New role of user account\n"+
		"    -lock\n"+
		"        Lock user\n"+
		"    -unlock\n"+
		"        Unlock user\n"+
		"    -enable\n"+
		"        Enable user\n"+
		"    -disable\n"+
		"        Disable user\n"+
		"    -password=<pass>\n"+
		"        New password. If omitted the password will be asked and read from stdin\n"+
		"    -password-file=<file>\n"+
		"        Read new password from <file>. The password MUST be the first line in the file, all other lines are ignored\n"+
		"\n", version)
}
