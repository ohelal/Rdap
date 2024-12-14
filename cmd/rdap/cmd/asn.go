package cmd

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/briandowns/spinner"
    "github.com/spf13/cobra"
)

var asnCmd = &cobra.Command{
    Use:   "asn [asn-number]",
    Short: "Query ASN information",
    Long: `Query RDAP information for an Autonomous System Number (ASN).

Examples:
  # Look up Google's ASN
  rdap asn 15169
  rdap asn AS15169

  # Look up Cloudflare's ASN
  rdap asn 13335
  rdap asn AS13335

Note: ASN numbers must be positive integers. You can optionally prefix them with 'AS'.`,
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            return fmt.Errorf("requires exactly one ASN number argument")
        }

        // Remove AS prefix if present for validation
        asnStr := strings.ToUpper(args[0])
        if strings.HasPrefix(asnStr, "AS") {
            asnStr = strings.TrimPrefix(asnStr, "AS")
        }

        // Validate ASN number
        asn, err := strconv.ParseUint(asnStr, 10, 32)
        if err != nil {
            return fmt.Errorf("invalid ASN format: %s. ASN must be a positive integer (e.g., 15169 or AS15169)", args[0])
        }

        if asn == 0 {
            return fmt.Errorf("ASN number cannot be zero")
        }

        return nil
    },
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
                fmt.Print("✓ Result cached\n")
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
            if isStatusError(err) {
                switch getStatusCode(err) {
                case 400:
                    return fmt.Errorf("invalid ASN format: %s. ASN must be a positive integer", asnNumber)
                case 404:
                    return fmt.Errorf("ASN %s not found in the RDAP database", asnNumber)
                case 429:
                    return fmt.Errorf("rate limit exceeded, please try again later")
                default:
                    return fmt.Errorf("server error: %w", err)
                }
            }
            return fmt.Errorf("querying ASN: %w", err)
        }

        // Cache the result
        cacheResult("asn:"+asnNumber, result, 1*time.Hour)

        // Render the result based on output style
        renderASNResult("ASN", asnNumber, result)
        return nil
    },
}

// [Rest of the rendering functions remain unchanged...]

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
			for _, field := range vcardFields {
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

func isStatusError(err error) bool {
    if err == nil {
        return false
    }
    _, ok := err.(interface {
        StatusCode() int
    })
    return ok
}

func getStatusCode(err error) int {
    if statusErr, ok := err.(interface {
        StatusCode() int
    }); ok {
        return statusErr.StatusCode()
    }
    return 0
}