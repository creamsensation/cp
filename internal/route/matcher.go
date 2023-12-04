package route

import (
	"fmt"
	"regexp"
	"strings"
)

func createMatcher(path string) *regexp.Regexp {
	items := make([]string, 0)
	for _, item := range strings.Split(path, "/") {
		if len(item) == 0 {
			continue
		}
		if item == "*" {
			items = append(items, regexUrlWildcard)
			continue
		}
		if strings.HasPrefix(item, "{") && strings.HasSuffix(item, "}") {
			item = strings.TrimSuffix(strings.TrimPrefix(item, "{"), "}")
			hasMatcher := strings.Contains(item, ":")
			if hasMatcher {
				matcher := item[strings.Index(item, ":")+1:]
				item = item[:strings.Index(item, ":")]
				items = append(items, matcher)
			}
			if !hasMatcher {
				items = append(items, regexUrlVar)
			}
			continue
		}
		items = append(items, item)
	}
	path = strings.Join(items, "/")
	rgx := regexp.MustCompile(fmt.Sprintf(regexUrlMatcher, path))
	return rgx
}
