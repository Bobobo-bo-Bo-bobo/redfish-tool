package main

import (
	"flag"
	"fmt"
	"os"
	"redfish"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	insecure := flag.Bool("insecure", false, "Skip SSL certificate verification")
	verbose := flag.Bool("verbose", false, "Verbose operation")
	ask := flag.Bool("ask", false, "Ask for password")
	user := flag.String("user", "", "Username to use for authentication")
	password := flag.String("password", "", "Password to use for authentication")
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
			Verbose:     *verbose,
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
