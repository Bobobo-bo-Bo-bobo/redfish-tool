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
		"Usage redfish-tool [-ask] [-config=<cfg>] [-help] [-password=<pass>] [-user=<user>] <command> [<cmd_options>] [<host> <host>]\n"+
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
		"\n", version)
}