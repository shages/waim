package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type colorFn func(a ...interface{}) string

// FormatPerms formats a file's 9-bit perms plus link/dir type
func FormatPerms(f *File, targetPerms PermType, color colorFn) string {
	u := f.userPerms()
	g := f.groupPerms()
	o := f.otherPerms()
	dir := "-"
	if f.IsLink() {
		dir = "l"
	} else if f.IsDir() {
		dir = "d"
	}
	if targetPerms == permUser {
		return fmt.Sprintf("%s%s%s%s", dir, color(u.format()), g.format(), o.format())
	} else if targetPerms == permGroup {
		return fmt.Sprintf("%s%s%s%s", dir, u.format(), color(g.format()), o.format())
	} else {
		return fmt.Sprintf("%s%s%s%s", dir, u.format(), g.format(), color(o.format()))
	}
}

// ListElaboratedPath traces an elaborated path, printing details of each path part
func ListElaboratedPath(elaborated ElaboratedPath, prefix string, usr User, action Action, pathCache map[string]bool) {
	if pathCache == nil {
		pathCache = (make(map[string]bool))
	}

	green := color.New(color.FgGreen).SprintFunc()
	blueBold := color.New(color.Bold, color.FgBlue).SprintFunc()
	redBold := color.New(color.Bold, color.FgRed).SprintFunc()

	lastPathCached := false
	for _, node := range elaborated.parts {
		// If we're disabling printing of duplicate paths, which happens when a linked
		// path is part of the same root tree as the original path, we cache the
		// paths encountered here
		if _, ok := pathCache[node.Path()]; ok {
			if !lastPathCached {
				fmt.Println(prefix + "...")
			}
			lastPathCached = true
			continue
		} else {
			lastPathCached = false
		}
		pathCache[node.Path()] = true

		switch n := node.(type) {
		case NodeUnknown:
			perms := "?????????? " + strings.Repeat("?", 10)
			fmt.Println(prefix + "? " + fmt.Sprintf("%s %s", perms, n.Path()))
		case NodeKnown:
			if n.target != nil {
				ListElaboratedPath(*n.target, prefix+green("| "), usr, action, pathCache)
				fmt.Println(prefix + green("|/"))
			}
			file := &n.file
			permType := usr.GetPermTypeOfFile(file)
			perms := file.getPermsByType(permType)

			canAccess := true
			if action == actionRead {
				if file.st.IsDir() && !perms.Execute {
					canAccess = false
				} else if !file.st.IsDir() && !perms.Read {
					canAccess = false
				}
			} else if action == actionExecute {
				if !file.st.IsDir() && !perms.Execute {
					canAccess = false
				}
			}

			statusStr := "* "
			if !canAccess {
				statusStr = redBold("X ")
			}
			prefixStr := prefix + statusStr
			permTargetColor := blueBold
			if !canAccess {
				permTargetColor = redBold
			}

			permStr := FormatPerms(file, permType, permTargetColor)
			groupStr := fmt.Sprintf("%-"+fmt.Sprintf("%d", elaborated.GroupMaxLen())+"v", file.Group.Name)
			ownerStr := fmt.Sprintf("%-"+fmt.Sprintf("%d", elaborated.OwnerMaxLen())+"v", file.User.Username)
			pathStr := n.Path()
			fmt.Printf("%s%s %s %s %s\n", prefixStr, permStr, groupStr, ownerStr, pathStr)
		}
	}
}
