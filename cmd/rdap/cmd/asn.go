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
        headers := []string{"Field", "Value"}
        var rows [][]string

        // Add basic information
        if mapData, ok := data.(map[string]interface{}); ok {
            if handle, ok := mapData["handle"].(string); ok {
                rows = append(rows, []string{"Name", handle})
            }
            if name, ok := mapData["name"].(string); ok {
                rows = append(rows, []string{"Name", name})
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

            // Add entity information
            if entities, ok := mapData["entities"].([]interface{}); ok {
                for _, entity := range entities {
                    if entityMap, ok := entity.(map[string]interface{}); ok {
                        if roles, ok := entityMap["roles"].([]interface{}); ok {
                            roleStr := make([]string, 0)
                            for _, r := range roles {
                                if str, ok := r.(string); ok {
                                    roleStr = append(roleStr, str)
                                }
                            }
                            if len(roleStr) > 0 {
                                rows = append(rows, []string{"Role", strings.Join(roleStr, ", ")})
                            }
                        }
                        if vcardArray, ok := entityMap["vcardArray"].([]interface{}); ok && len(vcardArray) > 1 {
                            if vcardData, ok := vcardArray[1].([]interface{}); ok {
                                for _, field := range vcardData {
                                    if fieldData, ok := field.([]interface{}); ok && len(fieldData) >= 3 {
                                        switch fieldData[0] {
                                        case "fn":
                                            if name, ok := fieldData[3].(string); ok {
                                                rows = append(rows, []string{"Contact Name", name})
                                            }
                                        case "email":
                                            if email, ok := fieldData[3].(string); ok {
                                                rows = append(rows, []string{"Email", email})
                                            }
                                        case "tel":
                                            if phone, ok := fieldData[3].(string); ok {
                                                rows = append(rows, []string{"Phone", phone})
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
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
    rootCmd.AddCommand(asnCmd)
}
