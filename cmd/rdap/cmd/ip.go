package cmd

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/briandowns/spinner"
    "github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
    Use:   "ip [ip-address]",
    Short: "Query IP information",
    Long: `Query RDAP information for an IP address.
Example: rdap ip 8.8.8.8`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        ipAddress := args[0]
        
        // Check cache first
        if cached, found := getCachedResult("ip:" + ipAddress); found {
            if data, ok := cached.(map[string]interface{}); ok {
                renderIPResult("IP", ipAddress, data)
                return nil
            }
        }
        
        // Show progress spinner if verbose
        var s *spinner.Spinner
        if verbose {
            s = newSpinner("Querying IP information...")
            s.Start()
            defer s.Stop()
        }

        // Query IP
        ctx := context.Background()
        result, err := client.QueryIP(ctx, ipAddress)
        if err != nil {
            return fmt.Errorf("querying IP: %w", err)
        }

        // Cache the result
        cacheResult("ip:"+ipAddress, result, 1*time.Hour)

        // Render the result based on output style
        renderIPResult("IP", ipAddress, result)
        return nil
    },
}

func renderIPResult(typ, query string, data interface{}) {
    switch outputStyle {
    case "table":
        headers := []string{"Field", "Value"}
        var rows [][]string

        // Add basic information
        if mapData, ok := data.(map[string]interface{}); ok {
            if handle, ok := mapData["handle"].(string); ok {
                rows = append(rows, []string{"Handle", handle})
            }
            if name, ok := mapData["name"].(string); ok {
                rows = append(rows, []string{"Name", name})
            }
            if ipVersion, ok := mapData["ipVersion"].(string); ok {
                rows = append(rows, []string{"IP Version", ipVersion})
            }
            if startAddr, ok := mapData["startAddress"].(string); ok {
                rows = append(rows, []string{"Start Address", startAddr})
            }
            if endAddr, ok := mapData["endAddress"].(string); ok {
                rows = append(rows, []string{"End Address", endAddr})
            }
            if status, ok := mapData["status"].([]interface{}); ok {
                statusStr := make([]string, 0)
                for _, s := range status {
                    if str, ok := s.(string); ok {
                        statusStr = append(statusStr, str)
                    }
                }
                if len(statusStr) > 0 {
                    rows = append(rows, []string{"Status", strings.Join(statusStr, ", ")})
                }
            }
        }

        renderTable(headers, rows)
    case "box":
        if mapData, ok := data.(map[string]interface{}); ok {
            renderBox(fmt.Sprintf("%s Query Result", typ), formatRDAPResult(mapData))
        }
    default:
        if mapData, ok := data.(map[string]interface{}); ok {
            if format == "json" {
                fmt.Println(formatJSON(mapData, true))
            } else {
                fmt.Print(formatRDAPResult(mapData))
            }
        }
    }
}

func init() {
    rootCmd.AddCommand(ipCmd)
}
