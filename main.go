package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

type myUser struct {
	*user.User
}

func (u *myUser) hasGroup(group string) (bool, error) {
	groups, err := u.GroupIds()
	if err != nil {
		return false, err
	}

	for _, b := range groups {
		if b == group {
			return true, nil
		}
	}
	return false, nil
}

func (u *myUser) canAccessPath(path string) (bool, string) {
	st, err := os.Stat(path)
	if os.IsPermission(err) {
		return false, ""
	} else if os.IsNotExist(err) {
		return false, "path does not exist"
	}

	file := myFile{st: st}
	hasGroup, _ := u.hasGroup(fmt.Sprint(file.Info().Gid))

	if fmt.Sprint(file.Info().Uid) == u.Uid {
		if !file.st.IsDir() {
			if !file.User().Read {
				return false, "file missing user read: " + path
			}
		} else {
			if !file.User().Execute {
				return false, "path missing user execute: " + path
			}
		}
	} else if hasGroup {
		if !file.st.IsDir() {
			if !file.Group().Read {
				return false, "file missing group read: " + path
			}
		} else {
			if !file.Group().Execute {
				return false, "path missing group execute: " + path
			}
		}
	} else {
		if !file.st.IsDir() {
			if !file.Other().Read {
				return false, "file missing other read: " + path
			}
		} else {
			if !file.Other().Execute {
				return false, "path missing other execute: " + path
			}
		}
	}
	return true, ""
}

type myFile struct {
	st os.FileInfo
}

func (f *myFile) Info() *syscall.Stat_t {
	return f.st.Sys().(*syscall.Stat_t)
}

type perm struct {
	Read    bool
	Write   bool
	Execute bool
}

func selectPerms(mode os.FileMode, offset uint8) perm {
	return perm{
		Read:    mode&(1<<(offset+2)) != 0,
		Write:   mode&(1<<(offset+1)) != 0,
		Execute: mode&(1<<offset) != 0,
	}
}

func (f *myFile) User() perm {
	return selectPerms(f.st.Mode(), 6)
}

func (f *myFile) Group() perm {
	return selectPerms(f.st.Mode(), 3)
}

func (f *myFile) Other() perm {
	return selectPerms(f.st.Mode(), 0)
}

func pathParts(path string) []string {
	parts := strings.Split(path, "/")[1:]
	paths := make([]string, 0)
	for i := range parts {
		str := "/" + strings.Join(parts[:i+1], "/")
		paths = append(paths, str)
	}
	return paths
}

func run() (int, error) {
	userPtr := flag.String("user", "__unknown", "user to check as, defaults to current user")
	flag.Parse()

	var usr myUser
	if *userPtr == "__unknown" {
		_usr, err := user.Current()
		if err != nil {
			return 1, err
		}
		usr = myUser{_usr}
	} else {
		userObj, err := user.Lookup(*userPtr)
		if err != nil {
			return 1, err
		}
		usr = myUser{userObj}
	}

        userPath := os.Args[len(os.Args)-1]
        absPath, err := filepath.Abs(userPath)

        if err != nil {
                return 1, err
        }

        for _, path := range pathParts(absPath) {
                canAccess, reason := usr.canAccessPath(path)
                if !canAccess {
                        fmt.Println("Cannot access path " + path + " because " + reason)
                        return 0, nil
                }
        }
        fmt.Println("... but you can access this path")
	return 0, nil
}

func main() {
	rc, err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(rc)
}
