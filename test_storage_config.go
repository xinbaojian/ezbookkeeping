package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

func main() {
	fmt.Println("Testing ezBookkeeping storage configuration validation...")

	// Test 1: Default configuration
	fmt.Println("\nTest 1: Default configuration (storage_path='storage/')")
	testConfig("D:\\workspace\\go\\src\\ezbookkeeping\\conf\\ezbookkeeping.ini")

	// Test 2: Path with forward slashes (absolute)
	fmt.Println("\nTest 2: Absolute path with forward slashes")
	createTestConfig("storage_type = local_filesystem\nlocal_filesystem_path = C:/tmp/ezbookkeeping_storage")

	// Test 3: Path with backslashes (absolute)
	fmt.Println("\nTest 3: Absolute path with backslashes")
	createTestConfig("storage_type = local_filesystem\nlocal_filesystem_path = C:\\tmp\\ezbookkeeping_storage")

	// Test 4: Path with environment variable
	fmt.Println("\nTest 4: Path with environment variable")
	os.Setenv("EBK_STORAGE_LOCAL_FILESYSTEM_PATH", "C:\\tmp\\ezbookkeeping_storage")
	createTestConfig("storage_type = local_filesystem\nlocal_filesystem_path = storage/")

	// Test 5: Relative path that doesn't exist
	fmt.Println("\nTest 5: Relative path that doesn't exist")
	createTestConfig("storage_type = local_filesystem\nlocal_filesystem_path = nonexistent/path")

	// Test 6: Alternative local_filesystem_path structure
	fmt.Println("\nTest 6: Alternative structure")
	createTestConfig("storage_type = local_filesystem\nlocal_filesystem_path = /tmp/ezbookkeeping_storage")

	os.Remove("test_config.ini")
}

func createTestConfig(configContent string) {
	err := os.WriteFile("test_config.ini", []byte("[storage]\n"+configContent), 0644)
	if err != nil {
		log.Fatal(err)
	}
	testConfig("test_config.ini")
}

func testConfig(configPath string) {
	cfg, err := settings.LoadConfiguration(configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		return
	}
	fmt.Printf("Local filesystem path: %s\n", cfg.LocalFileSystemPath)
	fmt.Printf("Storage type: %s\n", cfg.StorageType)

	// Check if path validation was successful
	if cfg.StorageType == settings.LocalFileSystemObjectStorageType {
		lastSep := filepath.Separator
		if cfg.LocalFileSystemPath[len(cfg.LocalFileSystemPath)-1] != byte(lastSep) {
			cfg.LocalFileSystemPath += string(lastSep)
		}
		if _, err := os.Stat(cfg.LocalFileSystemPath); err == nil {
			fmt.Println("✓ Path exists and is accessible")
		} else {
			fmt.Printf("✗ Path validation failed: %v\n", err)
		}
	}
}
