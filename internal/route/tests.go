package route

import cfg "github.com/creamsensation/cp/internal/config"

var (
	testsLangs = cfg.Languages{
		"cs": cfg.Language{Enabled: true, Default: true},
		"en": cfg.Language{Enabled: true},
	}
)
