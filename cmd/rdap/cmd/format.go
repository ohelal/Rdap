package cmd

import (
    "encoding/json"
    "fmt"
    "strings"
)

func formatRDAPResult(data map[string]interface{}) string {
    var sb strings.Builder

    // Format basic information
    if handle, ok := data["handle"].(string); ok {
        sb.WriteString(fmt.Sprintf("Handle: %s\n", handle))
    }
    if name, ok := data["name"].(string); ok {
        sb.WriteString(fmt.Sprintf("Name: %s\n", name))
    }
    
    // Format status
    if status, ok := data["status"].([]interface{}); ok {
        statusStr := make([]string, 0)
        for _, s := range status {
            if str, ok := s.(string); ok {
                statusStr = append(statusStr, str)
            }
        }
        if len(statusStr) > 0 {
            sb.WriteString(fmt.Sprintf("Status: %s\n", strings.Join(statusStr, ", ")))
        }
    }

    // Format entities
    if entities, ok := data["entities"].([]interface{}); ok {
        for _, e := range entities {
            if entity, ok := e.(map[string]interface{}); ok {
                sb.WriteString("\nEntity:\n")
                if roles, ok := entity["roles"].([]interface{}); ok {
                    roleStr := make([]string, 0)
                    for _, r := range roles {
                        if role, ok := r.(string); ok {
                            roleStr = append(roleStr, role)
                        }
                    }
                    if len(roleStr) > 0 {
                        sb.WriteString(fmt.Sprintf("  Roles: %s\n", strings.Join(roleStr, ", ")))
                    }
                }

                // Format vCard information
                if vcardArray, ok := entity["vcardArray"].([]interface{}); ok && len(vcardArray) > 1 {
                    if vcard, ok := vcardArray[1].([]interface{}); ok {
                        for _, v := range vcard {
                            if field, ok := v.([]interface{}); ok && len(field) >= 3 {
                                fieldName := field[0].(string)
                                if fieldValue, ok := field[3].(string); ok {
                                    switch fieldName {
                                    case "fn":
                                        sb.WriteString(fmt.Sprintf("  Name: %s\n", fieldValue))
                                    case "email":
                                        sb.WriteString(fmt.Sprintf("  Email: %s\n", fieldValue))
                                    case "tel":
                                        sb.WriteString(fmt.Sprintf("  Phone: %s\n", fieldValue))
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    return sb.String()
}

func formatJSON(data interface{}, pretty bool) string {
    var bytes []byte
    var err error
    if pretty {
        bytes, err = json.MarshalIndent(data, "", "  ")
    } else {
        bytes, err = json.Marshal(data)
    }
    if err != nil {
        return fmt.Sprintf("Error formatting JSON: %v", err)
    }
    return string(bytes)
}
