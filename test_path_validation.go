package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// Simulate the config loading process for storage section
func main() {
	fmt.Println("Testing path validation logic...")

	testCases := []struct {
		name string
		path string
	}{
		{"default", "storage/"},
		{"absolute forward slash", "C:/tmp/storage"},
		{"absolute backslash", "C:\\tmp\\storage"},
		{"relative with dot", "./storage"},
		{"relative parent", "../storage"},
	}

	workingPath, _ := os.Getwd()
	fmt.Printf("Working path: %s\n", workingPath)

	for _, tc := range testCases {
		fmt.Println("\n--- Testing:", tc.name, "---")
		fmt.Printf("Input path: %s\n", tc.path)

		// Simulate the getFinalPath function
		p := tc.path
		if !filepath.IsAbs(p) {
			p = filepath.Join(workingPath, p)
			fmt.Printf("Joined path: %s\n", p)
		}

		// Test os.Stat
		_, err := os.Stat(p)
		if err != nil {
			fmt.Printf("os.Stat error: %v\n", err)
		} else {
			fmt.Println("Path exists ✓")
		}

		// Test with INI file
		iniContent := fmt.Sprintf("[storage]\ntype = local_filesystem\nlocal_filesystem_path = %s", tc.path)
		cfg, _ := ini.LoadSources(ini.LoadOptions{
			IgnoreInlineComment: true,
		}, []byte(iniContent))

		parsedPath := cfg.Section("storage").Key("local_filesystem_path").String()
		fmt.Printf("Parsed from INI: %s\n", parsedPath)
	}

	// Test environment variable preprocessing
	fmt.Println("\n--- Testing Environment Variable Preprocessing ---")
	os.Setenv("EBK_STORAGE_TYPE", "local_filesystem")
	os.Setenv("EBK_STORAGE_LOCAL_FILESYSTEM_PATH", "C:\\tmp\\custom_storage")

	testContent := `
[storage]
type = local_filesystem
local_filesystem_path = default_storage
`
	cfg, _ := ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, []byte(testContent))

	// Simulate getConfigItemStringValue
	sectionName := "storage"
	itemName := "local_filesystem_path"
	environmentValue := os.Getenv(fmt.Sprintf("EBK_%s_%s", strings.ToUpper(sectionName), strings.ToUpper(itemName)))

	fmt.Printf("Environment value: %s\n", environmentValue)
	fmt.Printf("INI value: %s\n", cfg.Section(sectionName).Key(itemName).String())

	// Test the path resolution logic
	workingDir, _ := os.Getwd()
	if environmentValue != "" {
		parsedPath := environmentValue
		if !filepath.IsAbs(parsedPath) {
			parsedPath = filepath.Join(workingDir, parsedPath)
		}
		fmt.Printf("Final resolved path: %s\n", parsedPath)
	}
}
