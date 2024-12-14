package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func formatRDAPResult(data map[string]interface{}) string {
	var sb strings.Builder

	formatBasicInfo(&sb, data)
	formatStatus(&sb, data)
	formatEntities(&sb, data)
	formatEvents(&sb, data)
	formatRemarks(&sb, data)

	return sb.String()
}

func formatBasicInfo(sb *strings.Builder, data map[string]interface{}) {
	if handle, ok := data["handle"].(string); ok {
		sb.WriteString(fmt.Sprintf("Handle: %s\n", handle))
	}
	if name, ok := data["name"].(string); ok {
		sb.WriteString(fmt.Sprintf("Name: %s\n", name))
	}
}

func formatStatus(sb *strings.Builder, data map[string]interface{}) {
	if status, ok := data["status"].([]interface{}); ok {
		statusStr := getStringSlice(status)
		if len(statusStr) > 0 {
			sb.WriteString(fmt.Sprintf("Status: %s\n", strings.Join(statusStr, ", ")))
		}
	}
}

func formatEntities(sb *strings.Builder, data map[string]interface{}) {
	if entities, ok := data["entities"].([]interface{}); ok {
		for _, e := range entities {
			if entity, ok := e.(map[string]interface{}); ok {
				sb.WriteString("\nEntity:\n")
				formatEntityRoles(sb, entity)
				formatVCardInfo(sb, entity)
			}
		}
	}
}

func formatEntityRoles(sb *strings.Builder, entity map[string]interface{}) {
	if roles, ok := entity["roles"].([]interface{}); ok {
		roleStr := getStringSlice(roles)
		if len(roleStr) > 0 {
			sb.WriteString(fmt.Sprintf("  Roles: %s\n", strings.Join(roleStr, ", ")))
		}
	}
}

func formatVCardInfo(sb *strings.Builder, entity map[string]interface{}) {
	if vcardArray, ok := entity["vcardArray"].([]interface{}); ok && len(vcardArray) > 1 {
		if vcardData, ok := vcardArray[1].([]interface{}); ok {
			for _, field := range vcardData {
				formatVCardField(sb, field)
			}
		}
	}
}

func formatVCardField(sb *strings.Builder, field interface{}) {
	if fieldData, ok := field.([]interface{}); ok && len(fieldData) >= 3 {
		switch fieldData[0] {
		case "fn":
			if name, ok := fieldData[3].(string); ok {
				sb.WriteString(fmt.Sprintf("  Name: %s\n", name))
			}
		case "email":
			if email, ok := fieldData[3].(string); ok {
				sb.WriteString(fmt.Sprintf("  Email: %s\n", email))
			}
		case "tel":
			if phone, ok := fieldData[3].(string); ok {
				sb.WriteString(fmt.Sprintf("  Phone: %s\n", phone))
			}
		}
	}
}

func formatEvents(sb *strings.Builder, data map[string]interface{}) {
	if events, ok := data["events"].([]interface{}); ok {
		for _, event := range events {
			if eventMap, ok := event.(map[string]interface{}); ok {
				formatEventDetails(sb, eventMap)
			}
		}
	}
}

func formatEventDetails(sb *strings.Builder, event map[string]interface{}) {
	if action, ok := event["eventAction"].(string); ok {
		if date, ok := event["eventDate"].(string); ok {
			sb.WriteString(fmt.Sprintf("\nEvent (%s): %s\n", action, date))
		}
	}
}

func formatRemarks(sb *strings.Builder, data map[string]interface{}) {
	if remarks, ok := data["remarks"].([]interface{}); ok {
		for _, remark := range remarks {
			if remarkMap, ok := remark.(map[string]interface{}); ok {
				formatRemarkDetails(sb, remarkMap)
			}
		}
	}
}

func formatRemarkDetails(sb *strings.Builder, remark map[string]interface{}) {
	if description, ok := remark["description"].([]interface{}); ok {
		descStr := getStringSlice(description)
		if len(descStr) > 0 {
			sb.WriteString("\nRemark:\n")
			for _, desc := range descStr {
				sb.WriteString(fmt.Sprintf("  %s\n", desc))
			}
		}
	}
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

func getStringSlice(data []interface{}) []string {
	var strSlice []string
	for _, s := range data {
		if str, ok := s.(string); ok {
			strSlice = append(strSlice, str)
		}
	}
	return strSlice
}
