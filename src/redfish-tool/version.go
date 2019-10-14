package main

import (
	"fmt"
	redfish "git.ypbind.de/repository/go-redfish.git"
	"runtime"
)

func showVersion() {

	fmt.Printf("redfish-tool version %s\n"+
		"Copyright (C) by Andreas Maus <maus@ypbind.de>\n"+
		"This program comes with ABSOLUTELY NO WARRANTY.\n"+
		"\n"+
		"redfish-tool is distributed under the Terms of the GNU General\n"+
		"Public License Version 3. (http://www.gnu.org/copyleft/gpl.html)\n"+
		"\n"+
		"Build with go version: %s\n"+
		"Using go-redfish version: %s\n"+
		"\n", version, runtime.Version(), redfish.GoRedfishVersion)
}
