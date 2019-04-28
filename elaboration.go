package main

import (
	"os"
	"path/filepath"
	"strings"
)

func ElaboratePath(targetPath string) (ElaboratedPath, error) {
	parts := []Node{}
	groups := make(map[string]bool)
	owners := make(map[string]bool)
	for _, path := range pathParts(targetPath) {
		st, err := os.Lstat(path)
		if err != nil {
			// TODO: Handle permission vs not found errors
			parts = append(parts, NodeUnknown{path: path})
			continue
		}

		file, err := MakeFile(st)
		if err != nil {
			return ElaboratedPath{}, err
		}
		groups[file.Group.Name] = true
		owners[file.User.Username] = true

		if file.IsLink() {
			resolvedLink, _ := os.Readlink(path)
			if resolvedLink[0] != '/' {
				parentPath := filepath.Dir(path)
				comb := parentPath + "/" + resolvedLink
				resolvedLink = filepath.Clean(comb)
			}

			elaboratedTarget, err := ElaboratePath(resolvedLink)
			if err != nil {
				return ElaboratedPath{}, err
			}

			parts = append(parts, NodeKnown{
				path:   path,
				st:     st,
				file:   *file,
				target: &elaboratedTarget,
			})
		} else {
			parts = append(parts, NodeKnown{path: path, st: st, file: *file, target: nil})
		}
	}
	return ElaboratedPath{parts: parts, groups: groups, owners: owners}, nil
}

// Break a given path string into an array of strings
func pathParts(path string) []string {
	pathSep := string(filepath.Separator)
	parts := strings.Split(path, pathSep)[1:]
	paths := make([]string, 0)
	for i := range parts {
		str := pathSep + strings.Join(parts[:i+1], pathSep)
		paths = append(paths, str)
	}
	return paths
}
