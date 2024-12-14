package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage RDAP configuration",
    Long: `Manage the RDAP configuration. Available subcommands:
  - init: Initialize default configuration
  - view: View current configuration`,
}

var configInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize default configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := viper.SafeWriteConfig(); err != nil {
            fmt.Printf("%s\n", errorStyle("Failed to initialize config"))
            return err
        }
        fmt.Println(successStyle("Default configuration file created"))
        return nil
    },
}

var configViewCmd = &cobra.Command{
    Use:   "view",
    Short: "View current configuration",
    RunE: func(cmd *cobra.Command, args []string) error {
        settings := viper.AllSettings()
        switch outputStyle {
        case "json":
            fmt.Println(formatJSON(settings, true))
        case "table":
            headers := []string{"Setting", "Value"}
            var rows [][]string
            for k, v := range settings {
                rows = append(rows, []string{k, fmt.Sprintf("%v", v)})
            }
            renderTable(headers, rows)
        default:
            renderBox("Current Configuration", formatJSON(settings, true))
        }
        return nil
    },
}

func init() {
    rootCmd.AddCommand(configCmd)
    configCmd.AddCommand(configInitCmd)
    configCmd.AddCommand(configViewCmd)
}
