package vscode_extensions

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/osquery/osquery-go/plugin/table"
	"github.com/pkg/errors"
)

var homeDirLocations = map[string]string{
	"windows": "/Users",
	"darwin":  "/Users",
	"linux":   "/home/",
}

var extensionsDir = map[string][]string{
	"windows": {".vscode/extensions", ".vscode-server/extensions/"},
	"darwin":  {".vscode/extensions", ".vscode-server/extensions/"},
	"linux":   {".vscode/extensions", ".vscode-server/extensions/"},
}

type userFileInfo struct {
	user string
	path string
}

type Extension struct {
	Name        string     `json:"name"`
	DisplayName string     `json:"displayName"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Publisher   string     `json:"publisher"`
	License     string     `json:"license"`
	Repository  Repository `json:"repository"`
	Categories  []string   `json:"categories"`
	Metadata    Metadata   `json:"__metadata"`
}

type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Metadata struct {
	ID                   string `json:"id"`
	PublisherID          string `json:"publisherId"`
	PublisherDisplayName string `json:"publisherDisplayName"`
	InstalledTimestamp   int64  `json:"installedTimestamp"`
}

func VSCodeColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("name"),
		table.TextColumn("category"),
		table.TextColumn("description"),
		table.TextColumn("display_name"),
		table.TextColumn("license"),
		table.TextColumn("path"),
		table.TextColumn("url"),
		table.TextColumn("version"),
		table.TextColumn("extension_id"),
		table.TextColumn("identifier"),
		table.TextColumn("publisher"),
		table.TextColumn("publisher_id"),
		table.BigIntColumn("installed_at"),
		table.TextColumn("user"),
	}
}

func parseExtension(ctx context.Context, fileInfo userFileInfo) (map[string]string, error) {
	var results map[string]string
	data, err := os.ReadFile(fileInfo.path)
	if err != nil {
		return nil, errors.Wrap(err, "reading extension metadata file")
	}
	var extensionInfo Extension
	if err := json.Unmarshal(data, &extensionInfo); err != nil {
		return nil, errors.Wrap(err, "unmarshalling extension metadata file")
	}

	results = map[string]string{
		"name":         extensionInfo.Name,
		"category":     strings.Join(extensionInfo.Categories, " - "),
		"description":  extensionInfo.Description,
		"display_name": extensionInfo.DisplayName,
		"license":      extensionInfo.License,
		"path":         fileInfo.path,
		"url":          extensionInfo.Repository.URL,
		"version":      extensionInfo.Version,
		"identifier":   extensionInfo.Publisher + "." + extensionInfo.Name,
		"extension_id": extensionInfo.Metadata.ID,
		"publisher":    extensionInfo.Metadata.PublisherDisplayName,
		"publisher_id": extensionInfo.Metadata.PublisherID,
		"user":         fileInfo.user,
		"installed_at": strconv.FormatInt(extensionInfo.Metadata.InstalledTimestamp, 10),
	}

	return results, nil
}

// Per docs generator function has to return an array of map of strings
func VSCodeExtGenerate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	osExtensionsDir := extensionsDir[runtime.GOOS]
	var results []map[string]string
	for _, extDir := range osExtensionsDir {
		userFiles, err := findFileInUserDirs(extDir)
		if err != nil {
			return results, nil
		}
		for _, file := range userFiles {
			res, err := parseExtension(ctx, file)
			if err != nil {
				continue
			}
			results = append(results, res)
		}
	}

	return results, nil
}

func findFileInUserDirs(extensions_location string) ([]userFileInfo, error) {
	foundPaths := []userFileInfo{}

	homedirRoots, ok := homeDirLocations[runtime.GOOS]
	if !ok {
		return []userFileInfo{}, errors.New("No homedir location found for this GOOS")
	}

	userDirs, _ := os.ReadDir(homedirRoots)

	// For each user's dir, in this possibleHome, check!
	for _, ud := range userDirs {
		userPathPattern := filepath.Join(homedirRoots, ud.Name(), extensions_location, "*/package.json")
		fullPaths, err := filepath.Glob(userPathPattern)
		if err != nil {
			continue
		}
		for _, fullPath := range fullPaths {
			if stat, err := os.Stat(fullPath); err == nil && stat.Mode().IsRegular() {
				foundPaths = append(foundPaths, userFileInfo{
					user: ud.Name(),
					path: fullPath,
				})
			}
		}
	}

	return foundPaths, nil
}
