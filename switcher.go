package cp

import (
	"fmt"
	"strings"
	
	"github.com/creamsensation/cp/internal/constant/requestVar"
)

type Switcher interface {
	Lang(langCode string, overwrite ...Map) string
}

type switcher struct {
	*control
}

func (s switcher) Lang(langCode string, overwrite ...Map) string {
	var vars Map
	if !s.core.router.localized {
		return fmt.Sprintf("<%s:router is not localized>", langCode)
	}
	currentLang := s.control.Request().Lang()
	index := -1
	for lc, lr := range s.core.router.localizedRoutes {
		if lc != currentLang {
			continue
		}
		for i, r := range lr {
			if r.Matcher.MatchString(s.control.Request().Path()) {
				index = i
				break
			}
		}
	}
	if index == -1 {
		return fmt.Sprintf("[%s:localized route does not exist>", langCode)
	}
	langRoutes, ok := s.core.router.localizedRoutes[langCode]
	if !ok {
		return fmt.Sprintf("<%s:language does not exist>", langCode)
	}
	if len(langRoutes)-1 < index {
		return fmt.Sprintf("/%s", langCode)
	}
	match := langRoutes[index]
	path := match.Path
	if len(overwrite) > 0 {
		vars = overwrite[0]
	}
	for vk, vp := range match.VarsPlaceholders {
		if vk == requestVar.Lang {
			path = strings.Replace(path, vp, langCode, 1)
			continue
		}
		v, ok := vars[vk]
		if ok {
			path = strings.Replace(path, vp, fmt.Sprintf("%v", v), 1)
		}
		if !ok {
			path = strings.Replace(path, vp, Var[string](s.control, vk), 1)
		}
	}
	return path
}
