package main

import (
	"fmt"
	"os/user"
)

// User is a thin wrapper around os/user.User to test applicable perms of a file to the
// user
type User struct {
	*user.User
}

// GetPermTypeOfFile returns the PermType of a File applicable to this User
func (u *User) GetPermTypeOfFile(f *File) PermType {
	hasGroup, err := u.isMemberOfGroup(f.Group)
	if err != nil {
		hasGroup = false
	}

	if fmt.Sprint(f.Uid) == u.Uid {
		return permUser
	} else if hasGroup {
		return permGroup
	} else {
		return permOther
	}
}

func (u *User) isMemberOfGroup(group *user.Group) (bool, error) {
	groups, err := u.GroupIds()
	if err != nil {
		return false, err
	}

	for _, b := range groups {
		if b == group.Gid {
			return true, nil
		}
	}
	return false, nil
}
