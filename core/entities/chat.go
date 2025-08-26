package entities

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type ProjectFile struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Path     string        `json:"path"`
	Children []ProjectFile `json:"children,omitempty"`
}

var idCounter int = 1

func nextID() string {
	id := strconv.Itoa(idCounter)
	idCounter++
	return id
}

func ScanDirOneLevel(path string) ([]ProjectFile, error) {
	var files []ProjectFile

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		pf := ProjectFile{
			ID:   nextID(),
			Name: name,
			Path: filepath.Join(path, name),
		}

		if entry.IsDir() {
			pf.Type = "folder"
			pf.Children = nil
		} else {
			pf.Type = "file"
		}

		files = append(files, pf)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].Type == files[j].Type {
			return files[i].Name < files[j].Name
		}
		return files[i].Type == "folder"
	})

	return files, nil
}

func SearchFilesRecursive(root string, search string) ([]ProjectFile, error) {
	ignoreDirs := map[string]struct{}{
		"node_modules": {},
		"vendor":       {},
		".git":         {},
		".idea":        {},
		".vscode":      {},
	}

	var scanDir func(path string) ([]ProjectFile, error)
	searchLower := strings.ToLower(search)

	scanDir = func(path string) ([]ProjectFile, error) {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var files []ProjectFile
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			if entry.IsDir() {
				if _, ignored := ignoreDirs[name]; ignored {
					continue
				}
			}

			pf := ProjectFile{
				ID:   nextID(),
				Name: name,
				Path: filepath.Join(path, name),
			}

			if entry.IsDir() {
				pf.Type = "folder"
				children, err := scanDir(pf.Path)
				if err != nil {
					return nil, err
				}
				if len(children) > 0 {
					pf.Children = children
					files = append(files, pf)
				}
			} else {
				pf.Type = "file"
				if search != "" {
					match, _ := filepath.Match(search, name)
					if !match && !strings.Contains(strings.ToLower(name), searchLower) {
						continue
					}
				}
				files = append(files, pf)
			}
		}

		sort.Slice(files, func(i, j int) bool {
			if files[i].Type == files[j].Type {
				return files[i].Name < files[j].Name
			}
			return files[i].Type == "folder"
		})

		return files, nil
	}

	return scanDir(root)
}
