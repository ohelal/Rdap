package cmd

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

var domainCmd = &cobra.Command{
	Use:   "domain [domain-name]",
	Short: "Query domain information",
	Long: `Query RDAP information for a domain name.
Example: rdap domain example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domainName := args[0]

		// Check cache first
		if cached, found := getCachedResult("domain:" + domainName); found {
			if data, ok := cached.(map[string]interface{}); ok {
				renderDomainResult("Domain", domainName, data)
				return nil
			}
		}

		// Show progress spinner if verbose
		var s *spinner.Spinner
		if verbose {
			s = newSpinner("Querying domain information...")
			s.Start()
			defer s.Stop()
		}

		// Query domain
		ctx := context.Background()
		result, err := client.QueryDomain(ctx, domainName)
		if err != nil {
			return fmt.Errorf("querying domain: %w", err)
		}

		// Cache the result
		cacheResult("domain:"+domainName, result, 1*time.Hour)

		// Render the result based on output style
		renderDomainResult("Domain", domainName, result)
		return nil
	},
}

func renderDomainResult(typ, query string, data interface{}) {
	switch outputStyle {
	case "table":
		headers := []string{"Field", "Value"}
		var rows [][]string

		// Add basic information
		if mapData, ok := data.(map[string]interface{}); ok {
			if handle, ok := mapData["handle"].(string); ok {
				rows = append(rows, []string{"Handle", handle})
			}
			if name, ok := mapData["ldhName"].(string); ok {
				rows = append(rows, []string{"Domain Name", name})
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
	rootCmd.AddCommand(domainCmd)
}
