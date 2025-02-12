package json2scheme

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// JsonToScheme converts a JSON raw message to Scheme format
func JsonToScheme(data json.RawMessage) (string, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return mapToScheme(m), nil
}

// mapToScheme converts a nested map to Scheme format
func mapToScheme(m interface{}) string {
	switch v := m.(type) {
	case map[string]interface{}:
		// Sort keys for consistent output
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		pairs := make([]string, 0, len(v))
		for _, key := range keys {
			value := v[key]
			pairs = append(pairs, fmt.Sprintf("(%s %s)", key, mapToScheme(value)))
		}
		return fmt.Sprintf("(%s)", strings.Join(pairs, " "))
	case []interface{}:
		elements := make([]string, len(v))
		for i, elem := range v {
			elements[i] = mapToScheme(elem)
		}
		return fmt.Sprintf("(%s)", strings.Join(elements, " "))
	case string:
		return fmt.Sprintf("%q", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "#t"
		}
		return "#f"
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%v", v)
	}
}
