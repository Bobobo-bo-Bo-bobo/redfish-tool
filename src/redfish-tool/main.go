package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	redfish "git.ypbind.de/repository/go-redfish.git"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	insecure := flag.Bool("insecure", false, "Skip SSL certificate verification")
	debug := flag.Bool("debug", false, "Debug operation")
	ask := flag.Bool("ask", false, "Ask for password")
	user := flag.String("user", "", "Username to use for authentication")
	password := flag.String("password", "", "Password to use for authentication")
	password_file := flag.String("password-file", "", "Read password from file")
	config_file := flag.String("config", "", "Configuration file to use")
	help := flag.Bool("help", false, "Show help text")
	hosts := flag.String("host", "", "Hosts to work on")
	port := flag.Int("port", 0, "Alternate port to connect to")
	timeout := flag.Int64("timeout", 60, "Connection timeout in seconds")

	flag.Usage = ShowUsage
	flag.Parse()
	if *help {
		ShowUsage()
		os.Exit(0)
	}

	trailing := flag.Args()

	if *config_file != "" {
		// read and parse configuration file
	} else {
		if *ask {
			fmt.Print("Password: ")
			raw_pass, _ := terminal.ReadPassword(int(syscall.Stdin))
			pass := strings.TrimSpace(string(raw_pass))
			password = &pass
			fmt.Println()
		}
		if *password_file != "" {
			_passwd, err := ReadSingleLine(*password_file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Unable to read password from file: %s\n", err.Error())
				os.Exit(1)
			}
			password = &_passwd
		}
	}

	// get requested command
	if len(trailing) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command defined\n\n")
		ShowUsage()
		os.Exit(1)
	}
	command := strings.ToLower(trailing[0])

	if *hosts == "" {
		fmt.Fprintf(os.Stderr, "Error: No destination host given\n\n")
		ShowUsage()
		os.Exit(1)
	}

	if *timeout < 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid timeout %d; must be >= 0\n\n", *timeout)
		os.Exit(2)
	}

	host_list := strings.Split(*hosts, ",")
	for _, host := range host_list {
		rf := redfish.Redfish{
			Hostname:    host,
			Port:        *port,
			Username:    *user,
			Password:    *password,
			InsecureSSL: *insecure,
			Debug:       *debug,
			Timeout:     time.Duration(*timeout) * time.Second,
		}

		// XXX: optionally parse additional and command specific command lines
		if command == "get-all-users" {
			err = GetAllUsers(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-user" {
			err = GetUser(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-all-roles" {
			err = GetAllRoles(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-role" {
			err = GetRole(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-all-managers" {
			err = GetAllManagers(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-manager" {
			err = GetManager(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-all-systems" {
			err = GetAllSystems(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "get-system" {
			err = GetSystem(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "gen-csr" {
			err = GenCSR(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "fetch-csr" {
			err = FetchCSR(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "import-cert" {
			err = ImportCertificate(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "reset-sp" {
			err = ResetSP(rf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "add-user" {
			err = AddUser(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "del-user" {
			err = DelUser(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "modify-user" {
			err = ModifyUser(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "passwd" {
			err = Passwd(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else if command == "system-power" {
			err = SystemPower(rf, trailing[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: Unknown command %s\n\n", command)
			ShowUsage()
			os.Exit(1)
		}
	}

	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
