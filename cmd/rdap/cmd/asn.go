package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var asnCmd = &cobra.Command{
	Use:   "asn [asn-number]",
	Short: "Query ASN information",
	Long: `Query RDAP information for an Autonomous System Number (ASN).
Example: rdap asn 15169
         rdap asn AS15169`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asnNumber := args[0]

		// Format ASN number
		asnNumber = strings.ToUpper(asnNumber)
		if !strings.HasPrefix(asnNumber, "AS") {
			asnNumber = "AS" + asnNumber
		}

		// Check cache first
		if cached, found := getCachedResult("asn:" + asnNumber); found {
			if data, ok := cached.(map[string]interface{}); ok {
				renderASNResult("ASN", asnNumber, data)
				return nil
			}
		}

		// Show progress spinner if verbose
		var s *spinner.Spinner
		if verbose {
			s = newSpinner("Querying ASN information...")
			s.Start()
			defer s.Stop()
		}

		// Query ASN
		ctx := context.Background()
		result, err := client.QueryASN(ctx, asnNumber)
		if err != nil {
			return fmt.Errorf("querying ASN: %w", err)
		}

		// Cache the result
		cacheResult("asn:"+asnNumber, result, 1*time.Hour)

		// Render the result based on output style
		renderASNResult("ASN", asnNumber, result)
		return nil
	},
}

func renderASNResult(typ, query string, data interface{}) {
	switch outputStyle {
	case "table":
		renderASNTable(data)
	case "box":
		renderASNBox(data)
	default:
		renderASNDefault(data)
	}
}

func renderASNTable(data interface{}) {
	headers := []string{"Field", "Value"}
	rows := buildASNTableRows(data)
	renderTable(headers, rows)
}

func buildASNTableRows(data interface{}) [][]string {
	var rows [][]string
	mapData, ok := data.(map[string]interface{})
	if !ok {
		return rows
	}

	// Add basic information
	rows = appendBasicInfo(rows, mapData)
	
	// Add entity information
	rows = appendEntityInfo(rows, mapData)

	// Add events information
	rows = appendEventsInfo(rows, mapData)

	return rows
}

func appendBasicInfo(rows [][]string, mapData map[string]interface{}) [][]string {
	if handle, ok := mapData["handle"].(string); ok {
		rows = append(rows, []string{"Handle", handle})
	}
	if name, ok := mapData["name"].(string); ok {
		rows = append(rows, []string{"Name", name})
	}
	if status, ok := mapData["status"].([]interface{}); ok {
		statusStr := getStringSlice(status)
		if len(statusStr) > 0 {
			rows = append(rows, []string{"Status", strings.Join(statusStr, ", ")})
		}
	}
	return rows
}

func appendEntityInfo(rows [][]string, mapData map[string]interface{}) [][]string {
	if entities, ok := mapData["entities"].([]interface{}); ok {
		for _, entity := range entities {
			if entityMap, ok := entity.(map[string]interface{}); ok {
				rows = appendEntityRoles(rows, entityMap)
				rows = appendEntityContact(rows, entityMap)
			}
		}
	}
	return rows
}

func appendEntityRoles(rows [][]string, entityMap map[string]interface{}) [][]string {
	if roles, ok := entityMap["roles"].([]interface{}); ok {
		roleStr := getStringSlice(roles)
		if len(roleStr) > 0 {
			rows = append(rows, []string{"Role", strings.Join(roleStr, ", ")})
		}
	}
	return rows
}

func appendEntityContact(rows [][]string, entityMap map[string]interface{}) [][]string {
	if vcardArray, ok := entityMap["vcardArray"].([]interface{}); ok && len(vcardArray) > 1 {
		if vcardFields, ok := vcardArray[1].([]interface{}); ok {
			rows = appendVCardFields(rows, vcardFields)
		}
	}
	return rows
}

func appendEventsInfo(rows [][]string, mapData map[string]interface{}) [][]string {
	if events, ok := mapData["events"].([]interface{}); ok {
		for _, event := range events {
			if eventMap, ok := event.(map[string]interface{}); ok {
				rows = appendEventDetails(rows, eventMap)
			}
		}
	}
	return rows
}

func appendEventDetails(rows [][]string, eventMap map[string]interface{}) [][]string {
	if action, ok := eventMap["eventAction"].(string); ok {
		if date, ok := eventMap["eventDate"].(string); ok {
			rows = append(rows, []string{
				fmt.Sprintf("Event (%s)", action),
				date,
			})
		}
	}
	return rows
}

func getStringSlice(items []interface{}) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

func renderASNBox(data interface{}) {
	if mapData, ok := data.(map[string]interface{}); ok {
		renderBox(fmt.Sprintf("ASN Query Result"), formatRDAPResult(mapData))
	}
}

func renderASNDefault(data interface{}) {
	if mapData, ok := data.(map[string]interface{}); ok {
		if format == "json" {
			fmt.Println(formatJSON(mapData, true))
		} else {
			fmt.Print(formatRDAPResult(mapData))
		}
	}
}

func init() {
	rootCmd.AddCommand(asnCmd)
}
