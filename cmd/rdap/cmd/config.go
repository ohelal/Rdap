package cmd

import (
    "fmt"
    "os"
    "path/filepath"
)

var defaultConfig = map[string]interface{}{
    "timeout":  "10s",
    "format":   "pretty",
    "style":    "default",
    "verbose":  false,
    "base_url": "",
}

func getConfigPath() string {
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("%s\n", errorStyle("Failed to get home directory"))
        return ""
    }
    return filepath.Join(home, ".rdap", "config.yaml")
}

func initConfig() error {
    configPath := getConfigPath()
    if configPath == "" {
        return fmt.Errorf("failed to get config path")
    }

    // Create config directory if it doesn't exist
    configDir := filepath.Dir(configPath)
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }

    // Check if config file already exists
    if _, err := os.Stat(configPath); err == nil {
        fmt.Printf("%s\n", successStyle("Config file already exists"))
        return nil
    }

    // Create default config file
    file, err := os.Create(configPath)
    if err != nil {
        return fmt.Errorf("failed to create config file: %w", err)
    }
    defer file.Close()

    // Write default config
    for key, value := range defaultConfig {
        fmt.Fprintf(file, "%s: %v\n", key, value)
    }

    fmt.Printf("%s\n", successStyle("Config file created"))
    return nil
}

func loadConfig() error {
    configPath := getConfigPath()
    if configPath == "" {
        return fmt.Errorf("failed to get config path")
    }

    // Check if config file exists
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        fmt.Printf("%s\n", errorStyle("Config file not found"))
        return nil
    }

    // Read config file
    data, err := os.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }

    fmt.Printf("%s\n", successStyle("Config loaded"))
    fmt.Println(string(data))
    return nil
}
