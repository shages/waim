package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

// Action refers to the type of action to test
type Action int

const (
	actionRead Action = iota
	actionWrite
	actionExecute
)

// PermType refers to one of the 3 x 3-bit file perms
type PermType int

const (
	permUser PermType = iota
	permGroup
	permOther
)

// File contains a real file's lstat results, plus relevant gid/uid details
type File struct {
	st    os.FileInfo
	Gid   string
	Group *user.Group
	Uid   string
	User  *user.User
}

// IsLink returns whether the File is a symbolic link or not
func (f *File) IsLink() bool {
	return f.st.Mode()&os.ModeSymlink == os.ModeSymlink
}

// IsDir returns whether the File is a directory or not
func (f *File) IsDir() bool {
	return f.st.IsDir()
}

// MakeFile instantiates a File struct, bubbling up any lookup errors
func MakeFile(st os.FileInfo) (*File, error) {
	stat := st.Sys().(*syscall.Stat_t)
	gid := fmt.Sprint(stat.Gid)
	group, err := user.LookupGroupId(gid)
	if err != nil {
		return nil, err
	}
	uid := fmt.Sprint(stat.Uid)
	user, err := user.LookupId(uid)
	if err != nil {
		return nil, err
	}
	return &File{st: st, Gid: gid, Group: group, Uid: uid, User: user}, nil
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

func (f *File) userPerms() perm {
	return selectPerms(f.st.Mode(), 6)
}

func (f *File) groupPerms() perm {
	return selectPerms(f.st.Mode(), 3)
}

func (f *File) otherPerms() perm {
	return selectPerms(f.st.Mode(), 0)
}

func (f *File) getPermsByType(ptype PermType) perm {
	if ptype == permUser {
		return f.userPerms()
	} else if ptype == permGroup {
		return f.groupPerms()
	} else {
		return f.otherPerms()
	}
}

const (
	bitexec = 1 << iota
	bitwrite
	bitread
)

func (p *perm) format() string {
	read := "-"
	write := "-"
	execute := "-"

	if p.Read {
		read = "r"
	}
	if p.Write {
		write = "w"
	}
	if p.Execute {
		execute = "x"
	}
	return fmt.Sprintf("%s%s%s", read, write, execute)
}
