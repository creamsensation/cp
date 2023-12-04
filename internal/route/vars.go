package route

import (
	"strings"
)

func createVarsPlaceholders(path string) map[string]string {
	result := make(map[string]string)
	for _, item := range strings.Split(path, "/") {
		if !strings.HasPrefix(item, "{") && !strings.HasSuffix(item, "}") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(item, "{"), "}")
		hasMatcher := strings.Contains(name, ":")
		if hasMatcher {
			name = name[:strings.Index(name, ":")]
		}
		result[name] = item
	}
	return result
}
