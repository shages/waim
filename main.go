package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Flags struct {
	read       bool
	exec       bool
	targetPath string
	user       User
}

func parseFlags() (Flags, error) {
	userPtr := flag.String("user", "", "User to test file permissions against (default is current user)")
	flagRead := flag.Bool("read", false, "Test read access to target path (default behavior). Mutually exclusive with -exec")
	flagExecute := flag.Bool("exec", false, "Test execute access to target path. Mutually exclusive with -read")
	flag.Parse()

	if *flagRead && *flagExecute {
		return Flags{}, errors.New("--read and --exec are mutually exclusive")
	}

	var usr User
	if *userPtr == "" {
		_usr, err := user.Current()
		if err != nil {
			return Flags{}, err
		}
		usr = User{_usr}
	} else {
		userObj, err := user.Lookup(*userPtr)
		if err != nil {
			return Flags{}, err
		}
		usr = User{userObj}
	}

	userPath := os.Args[len(os.Args)-1]
	absPath, err := filepath.Abs(userPath)
	if err != nil {
		return Flags{}, err
	}

	return Flags{
		read:       *flagRead,
		exec:       *flagExecute,
		targetPath: absPath,
		user:       usr,
	}, nil
}

func run() (int, error) {
	flags, err := parseFlags()
	if err != nil {
		return 1, err
	}

	action := actionRead
	if flags.exec {
		action = actionExecute
	}

	elaboratedPath, err := ElaboratePath(flags.targetPath)
	if err != nil {
		return 1, err
	}
	ListElaboratedPath(elaboratedPath, "", flags.user, action, nil)
	return 0, nil
}

func main() {
	rc, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(rc)
}
