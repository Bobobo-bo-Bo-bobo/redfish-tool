package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	redfish "git.ypbind.de/repository/go-redfish.git"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	var format = OutputText

	insecure := flag.Bool("insecure", false, "Skip SSL certificate verification")
	debug := flag.Bool("debug", false, "Debug operation")
	ask := flag.Bool("ask", false, "Ask for password")
	user := flag.String("user", "", "Username to use for authentication")
	password := flag.String("password", "", "Password to use for authentication")
	passwordFile := flag.String("password-file", "", "Read password from file")
	configFile := flag.String("config", "", "Configuration file to use")
	help := flag.Bool("help", false, "Show help text")
	hosts := flag.String("host", "", "Hosts to work on")
	port := flag.Int("port", 0, "Alternate port to connect to")
	timeout := flag.Int64("timeout", 60, "Connection timeout in seconds")
	verbose := flag.Bool("verbose", false, "Verbose operation")
	version := flag.Bool("version", false, "Show version")
	outFormat := flag.String("format", "text", "Output format (text, JSON)")

	// Logging setup
	var logFmt = new(log.TextFormatter)
	logFmt.FullTimestamp = true
	logFmt.TimestampFormat = time.RFC3339
	log.SetFormatter(logFmt)

	flag.Usage = showUsage
	flag.Parse()
	if *help {
		showUsage()
		os.Exit(0)
	}

	if *version {
		showVersion()
		os.Exit(0)
	}

	trailing := flag.Args()

	if *configFile != "" {
		// TODO: not implemented yet - read and parse configuration file
	} else {
		if *ask {
			fmt.Print("Password: ")
			rawPass, _ := terminal.ReadPassword(int(syscall.Stdin))
			pass := strings.Replace(strings.Replace(strings.Replace(string(rawPass), "\r", "", -1), "\n", "", -1), "\t", "", -1)
			password = &pass
			fmt.Println()
		}
		if *passwordFile != "" {
			_passwd, err := readSingleLine(*passwordFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Unable to read password from file: %s\n", err.Error())
				os.Exit(1)
			}
			password = &_passwd
		}
	}

	_format := strings.ToLower(strings.TrimSpace(*outFormat))
	if _format == "text" {
		format = OutputText
	} else if _format == "json" {
		format = OutputJSON
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid output format\n\n")
		showUsage()
		os.Exit(1)
	}

	// get requested command
	if len(trailing) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No command defined\n\n")
		showUsage()
		os.Exit(1)
	}
	command := strings.ToLower(trailing[0])

	if *hosts == "" {
		fmt.Fprintf(os.Stderr, "Error: No destination host given\n\n")
		showUsage()
		os.Exit(1)
	}

	if *user == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing login credentials (username and/or password)")
		showUsage()
		os.Exit(1)
	}

	if *timeout < 0 {
		fmt.Fprintf(os.Stderr, "Error: Invalid timeout %d; must be >= 0\n\n", *timeout)
		os.Exit(2)
	}

	hostList := strings.Split(*hosts, ",")
	for _, host := range hostList {
		if *verbose {
			log.WithFields(log.Fields{
				"hostname": host,
			}).Info("Connecting to host")
		}

		rf := redfish.Redfish{
			Hostname:    host,
			Port:        *port,
			Username:    *user,
			Password:    *password,
			InsecureSSL: *insecure,
			Debug:       *debug,
			Timeout:     time.Duration(*timeout) * time.Second,
			Verbose:     *verbose,
		}

		if command == "get-all-users" {
			err = getAllUsers(rf, format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-user" {
			err = getUser(rf, trailing[1:], format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-all-roles" {
			err = getAllRoles(rf, format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-role" {
			err = getRole(rf, trailing[1:], format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-all-managers" {
			err = getAllManagers(rf, format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-manager" {
			err = getManager(rf, trailing[1:], format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-all-systems" {
			err = getAllSystems(rf, format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-system" {
			err = getSystem(rf, trailing[1:], format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "gen-csr" {
			err = genCSR(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "fetch-csr" {
			err = fetchCSR(rf)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "import-cert" {
			err = importCertificate(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "reset-sp" {
			err = resetSP(rf)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "add-user" {
			err = addUser(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "del-user" {
			err = delUser(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "modify-user" {
			err = modifyUser(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "passwd" {
			err = passwd(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "system-power" {
			err = systemPower(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "get-license" {
			err = getLicense(rf, trailing[1:], format)
			if err != nil {
				log.Error(err.Error())
			}
		} else if command == "add-license" {
			err = addLicense(rf, trailing[1:])
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			log.WithFields(log.Fields{
				"command": command,
			}).Error("Unknown command")
			showUsage()
			os.Exit(1)
		}
	}

	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
